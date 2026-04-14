package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	clientoficial "github.com/ecsavigne/client_wa_oficial/v2/client"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	generalpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	"github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	responseig "github.com/ecsavigne/client_wa_oficial/v2/types/ig/response"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type QueryData map[string]any

func NewQueryData() QueryData {
	return make(map[string]any)
}

func (q QueryData) String() string {
	query := ""
	for k, v := range q {
		query += fmt.Sprintf("%s=%v&", k, v)
	}
	query = strings.TrimSuffix(query, "&")
	return query
}

func (q *QueryData) SetValue(key string, value any) {
	(*q)[key] = value
}

func (q QueryData) GetValue(key string) any {
	return q[key]
}

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

func DefaultRequest(config clientoficial.ConfigClient, method string, ePoint string, data any) (*http.Request, error) {
	var (
		req    *http.Request
		err    error
		buff   = &bytes.Buffer{}
		byt    = make([]byte, 0)
		urlStr = ""
	)

	switch config.GetType() {
	case clientoficial.TYPE_CONFIG_IG:
		ePoint = strings.TrimPrefix(ePoint, "/")
		urlStr = fmt.Sprintf("%s/%s/%s/%s", config.GetBaseUrl(), config.GetVersion(), config.GetUserID(), ePoint)
	// case clientoficial.TYPE_CONFIG_WPP:
	// 	urlStr = fmt.Sprintf("%s/%s", urlStr, ePoint)
	default:
		return nil, fmt.Errorf("unsupported config type: %s", config.GetType())
	}

	switch method {
	case http.MethodGet, http.MethodDelete:
		if queryData, ok := data.(QueryData); ok {
			ePoint = fmt.Sprintf("%s?%s", urlStr, queryData.String())
		}
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		switch v := data.(type) {
		case proto.Message:
			byt, err = protojson.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal proto message: %v", err)
			}
			buff.Write(byt)
		case map[string]any:
			byt, err = json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map data: %v", err)
			}
			buff.Write(byt)
		default:
			return nil, fmt.Errorf("unsupported data type: %T", v)
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

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

func Do(cl clientoficial.Client, req *http.Request, responseType response.ResponseType) response.Responser {
	clientHttp := createClientHttp2()
	errorMsg := &generalpbv1.ResponseError{}

	res, err := clientHttp.Do(req)

	if err != nil {
		errorMsg.SetType(types.TypeErrorInRequest)
		errorMsg.SetCode(types.CodeErrorInRequest)
		errorMsg.SetMessage(fmt.Sprintf("Type: %s. Error is: %s", types.MsgErrorInRequest, err.Error()))
		return response.NewResponse(errorMsg, response.ResponseError)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {

	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 400:
		log := fmt.Sprintf("Error in function Do. Code: %d, Status: %s, MetaError: %s, BodyError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"), string(b))
		errorMsg.SetType(types.TypeErrorBadRequest)
		errorMsg.SetCode(types.CodeErrorBadRequest)
		errorMsg.SetMessage(log)
		return response.NewResponse(errorMsg, response.ResponseError)
	case 401:
		log := fmt.Sprintf("Error in function Do. Code: %d, Message: %s, MetaError: %s, BodyError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"), string(b))
		errorMsg.SetType(types.TypeErrorUnauthorized)
		errorMsg.SetCode(types.CodeErrorUnauthorized)
		errorMsg.SetMessage(log)
		return response.NewResponse(errorMsg, response.ResponseError)
	case 404:
		log := fmt.Sprintf("Error in function Do. Code: %d, Message: %s, MetaError: %s, BodyError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"), string(b))
		errorMsg.SetType(types.TypeErrorUrlNotFound)
		errorMsg.SetCode(types.CodeErrorUrlNotFound)
		errorMsg.SetMessage(log)
		return response.NewResponse(errorMsg, response.ResponseError)
	case 200:
		switch cl.GetType() {
		case clientoficial.CLIENT_IG.String():
			return responseig.WrapperResponseRequest(b, responseType)
		// case clientoficial.CLIENT_WHATSAPP.String():
		// 	return response.JsonWrapperResponseRequest(b, responseType)
		default:
			return nil
		}
	default:
		log := fmt.Sprintf("Error in function Do. Code: %d, Message: %s, MetaError: %s, BodyError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"), string(b))
		errorMsg.SetType(types.TypeErrorInRequest)
		errorMsg.SetCode(types.CodeErrorInRequest)
		errorMsg.SetMessage(log)
		return response.NewResponse(errorMsg, response.ResponseError)
	}
}
