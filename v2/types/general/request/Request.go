package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	clientoficial "github.com/ecsavigne/client_wa_oficial/v2/client"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	generalpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	generalresponse "github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	"github.com/ecsavigne/client_wa_oficial/v2/types/internal"

	"golang.org/x/net/http2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// defaultHeader sets the Authorization and Content-Type headers of the given http.Request.
// It uses the configClient to get the token, and the first contentType if provided.
// If no contentType is provided, it defaults to "application/json".
func defaultHeader(c clientoficial.ConfigClient, request *http.Request, contentType ...string) {
	cT := "application/json"
	if len(contentType) > 0 {
		cT = contentType[0]
	}

	h := make(http.Header)
	h.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetToken()))
	h.Set("Content-Type", cT)
	request.Header = h
}

func multiparHeader(c clientoficial.ConfigClient, request *http.Request, contentType string) {
	h := make(http.Header)
	h.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetToken()))
	h.Set("Content-Type", contentType)
	request.Header = h
}

func DefaultRequest(config clientoficial.ConfigClient, method string, ePoint string, data any, isID ...bool) (*http.Request, error) {
	var (
		req    *http.Request
		err    error
		buff   = &bytes.Buffer{}
		byt    = make([]byte, 0)
		urlStr = ""
		_isID  = true
	)

	if len(isID) > 0 {
		_isID = isID[0]
	}

	switch config.GetType() {
	case clientoficial.TYPE_CONFIG_IG:
		ePoint = strings.TrimPrefix(ePoint, "/")
		if !_isID {
			urlStr, err = url.JoinPath(config.GetBaseUrl(), config.GetVersion(), ePoint)
		} else {
			urlStr, err = url.JoinPath(config.GetBaseUrl(), config.GetVersion(), config.GetUserID(), ePoint)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to join URL path: %v", err)
		}
		// case clientoficial.TYPE_CONFIG_WPP:
	default:
		return nil, fmt.Errorf("unsupported config type: %s", config.GetType())
	}

	switch method {
	case http.MethodGet:
		if queryData, ok := data.(types.QueryData); ok {
			if len(queryData) > 0 {
				urlStr = fmt.Sprintf("%s?%s", urlStr, queryData.String())
			} else {
				urlStr = fmt.Sprintf("%s", urlStr)
			}

		}
	case http.MethodDelete:
		byt, err = json.Marshal(data)
		if err == nil {
			switch data.(type) {
			case types.QueryData, map[string]any, nil:
				if len(byt) > 0 {
					buff.Write(byt)
				} else {
					buff.Write(nil)
				}
			default:
				return nil, fmt.Errorf("unsupported data type for DELETE method: %T", data)
			}
		}
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		switch v := data.(type) {
		case proto.Message:
			byt, err = protojson.MarshalOptions{UseProtoNames: true}.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal proto message: %v", err)
			}
			byt = internal.CleanDataEmpty(byt)
			buff.Write(byt)
		case map[string]any:
			byt, err = json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map data: %v", err)
			}
			buff.Write(byt)
		case nil:
			buff.Write(nil)
		default:
			return nil, fmt.Errorf("unsupported data type: %T", v)
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	urlStr = strings.ReplaceAll(urlStr, "%3F", "?")
	fmt.Println("EndPoint URL:\n", method, "  ", urlStr, " Data:\n", buff.String())

	req, err = http.NewRequestWithContext(context.Background(), method, urlStr, buff)
	defaultHeader(config, req)

	return req, err
}

func createClientHttp2(timeOut ...int) *http.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()

	tr.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
	http2.ConfigureTransport(tr)

	time_Out := 60
	if len(timeOut) > 0 {
		time_Out = timeOut[0]
	}
	d := time.Duration(time_Out)

	return &http.Client{
		Timeout:   d * time.Second,
		Transport: tr,
	}
}

// multipartRequest(methoth, ePoint, c.Config, msg)
// func MultipartRequest(config clientoficial.ConfigClient, method string, ePoint string, data proto.Message) (*http.Request, error) {
// 	var (
// 		e error
// 	)
// 	if !strings.HasPrefix(ePoint, "/") {
// 		e = fmt.Errorf("Error in multipartRequest, file: RequestConfig.go.Error is: EndPoint is not start with /")
// 		return nil, e
// 	}

// 	switch v := data.(type) {
// 	case *igpbv1.InstagramMediaMessage:
// 		return nil, nil
// 		// v.GetFileHeader()
// 		// return utilMultipartRequest(multipartRequestData{
// 		// 	config:           config,
// 		// 	method:           method,
// 		// 	ePoint:           ePoint,
// 		// 	fileInfo:         v.GetFileHeader(),
// 		// 	link:             v.GetMessage().GetAttachment().GetPayload().GetUrl(),
// 		// 	typeClient:       clientoficial.CLIENT_IG,
// 		// 	msgType:          v.GetMessage().GetAttachment().GetType(),
// 		// 	messagingProduct: "instagram",
// 		// })
// 	case *igpbv1.InstagramMediaShareMessage:
// 		return utilMultipartRequest(multipartRequestData{})
// 	default:
// 		e = fmt.Errorf("Error in multipartRequest, file: RequestConfig.go.Error is: Unsupported message type: %T", v)
// 		return nil, e
// 	}
// }

func Do(cl clientoficial.Client, req *http.Request, responseType generalresponse.ResponseType) generalresponse.Responser {
	clientHttp := createClientHttp2()
	errorMsg := &generalpbv1.ResponseError{}

	res, err := clientHttp.Do(req)

	if err != nil {
		errorMsg.SetType(types.TypeErrorInRequest)
		errorMsg.SetCode(types.CodeErrorInRequest)
		errorMsg.SetMessage(fmt.Sprintf("Type: %s. Error is: %s", types.MsgErrorInRequest, err.Error()))
		return generalresponse.NewResponse(errorMsg, generalresponse.ResponseError)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {

	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 400:
		errorMsg.SetType(types.TypeErrorBadRequest)
		errorMsg.SetCode(types.CodeErrorBadRequest)
		errorMsg.SetMessage(string(b))
		return generalresponse.NewResponse(errorMsg, generalresponse.ResponseError)
	case 401:
		errorMsg.SetType(types.TypeErrorUnauthorized)
		errorMsg.SetCode(types.CodeErrorUnauthorized)
		errorMsg.SetMessage(string(b))
		return generalresponse.NewResponse(errorMsg, generalresponse.ResponseError)
	case 403:
		errorMsg.SetType(types.TypeErrorForbidden)
		errorMsg.SetCode(types.CodeErrorForbidden)
		errorMsg.SetMessage(string(b))
		return generalresponse.NewResponse(errorMsg, generalresponse.ResponseError)
	case 404:
		errorMsg.SetType(types.TypeErrorUrlNotFound)
		errorMsg.SetCode(types.CodeErrorUrlNotFound)
		errorMsg.SetMessage(string(b))
		return generalresponse.NewResponse(errorMsg, generalresponse.ResponseError)
	case 200:
		switch cl.GetType() {
		case clientoficial.CLIENT_IG.String():
			return generalresponse.WrapperResponseRequest(b, generalresponse.WRAPPER_RESPONSE_IG, responseType)
		// case clientoficial.CLIENT_WHATSAPP.String():
		// 	return generalresponse.JsonWrapperResponseRequest(b, responseType)
		default:
			return nil
		}
	default:
		errorMsg.SetType(types.TypeErrorInRequest)
		errorMsg.SetCode(types.CodeErrorInRequest)
		errorMsg.SetMessage(string(b))
		return generalresponse.NewResponse(errorMsg, generalresponse.ResponseError)
	}
}
