//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path"
	"strings"

	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"golang.org/x/net/http2"
)

type TypeRequest = string

const (
	RequestGetMessageInfo TypeRequest = "RequestGetMessageInfo"
	RequestChangeUrlFull  TypeRequest = "RequestChangeUrlFull"
)

func defaultHeader(c *Config) {
	c.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	c.request.Header.Set("Content-Type", "application/json")
}

func multiparHeader(c *Config, contentType string) {
	c.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	c.request.Header.Set("Content-Type", contentType)
}

func resetError(c *Config) {
	if c.Error != nil {
		c.Error = nil
	}
}

// defaultRequest func
// methoth: GET, POST, PUT, PATCH, DELETE
// ePoint: EndPoint
// c: Config
// params: [Optional]
// - map[string]any. body request,dataForm.
// - types.Messager protocol of message to send,
// - name of request ej> "GetMessageInfo"
func defaultRequest(methoth string, ePoint string, c *Config, params ...any) (*http.Request, error) {
	resetError(c)
	var (
		e              error
		urlPath        *url.URL
		urlAlternative string = ""
		msg            types.Messager
		formData       *bytes.Buffer = bytes.NewBuffer([]byte{})
	)

	if len(params) > 0 && len(params) == 1 {
		switch v := params[0].(type) {
		case types.Messager:
			msg = v
		case map[string]any:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, &types.Error{
					Type:    types.TypeErrorUnrecognized,
					Code:    types.CodeErrorUnrecognized,
					Message: err.Error(),
				}
			}
			formData = bytes.NewBuffer(b)
		case TypeRequest:
			switch v {
			case RequestGetMessageInfo:
				c.path = ""
				c.BaseUrl.Path = ePoint
			case RequestChangeUrlFull:
				urlAlternative = strings.TrimPrefix(ePoint, "/")
			}
		}
	} else {
		// validation posible
	}

	if !strings.HasPrefix(ePoint, "/") && urlAlternative == "" {
		var log = "Error in deafultRequest, file: RequestConfig.go.Error is: EndPoint is not start with /"
		c.Error = fmt.Errorf("%s", log)
		panic(c.Error)
	}

	urlPath, _ = url.Parse(c.path + ePoint)
	urlPath = c.BaseUrl.ResolveReference(urlPath)
	if msg != nil {
		c.request, e = http.NewRequest(methoth, urlPath.String(), msg.ToJSONReader())
	} else {
		if urlAlternative != "" {
			c.request, e = http.NewRequest(methoth, urlAlternative, formData)
		} else {
			c.request, e = http.NewRequest(methoth, urlPath.String(), formData)
		}
	}
	if e != nil {
		c.Error = fmt.Errorf("Error in defaultRequest, NewRequest: %s. Error is: %s", c.BaseUrl, e.Error())
		return nil, c.Error
	}
	defaultHeader(c)
	return c.request, nil
}

func multipartRequest(methoth string, ePoint string, c *Config, msg types.Messager) (*http.Request, error) {
	resetError(c)
	var e error
	var urlPath *url.URL
	if !strings.HasPrefix(ePoint, "/") {
		var log = "Error in multipartRequest, file: RequestConfig.go.Error is: EndPoint is not start with /"
		c.Error = fmt.Errorf("%s", log)
		panic(c.Error)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12},
	}
	http2.ConfigureTransport(tr)
	client := &http.Client{Transport: tr}

	resp, err := client.Get(msg.GetMessageLink())
	if err != nil {
		c.Error = fmt.Errorf("Error in multiparRequest getting file. Error is: %s", err.Error())
		return nil, c.Error
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 400:
		log := fmt.Sprintf("Error in function makeRequest bad request of ClientWA. Message type: %s. error is: %s", msg.GetType(), resp.Status)
		c.Error = &types.Error{
			Type:    types.TypeErrorBadRequest,
			Code:    types.CodeErrorBadRequest,
			Message: log,
		}
		return nil, c.Error
	case 401:
		log := fmt.Sprintf("Error in function makeRequest bad request of ClientWA. Message type: %s. error is: %s", msg.GetType(), resp.Status)
		c.Error = &types.Error{
			Type:    types.TypeErrorUnauthorized,
			Code:    types.CodeErrorUnauthorized,
			Message: log,
		}
		return nil, c.Error
	case 404:
		log := fmt.Sprintf("Error in function makeRequest bad request of ClientWA. Message type: %s. error is: %s", msg.GetType(), resp.Status)
		c.Error = &types.Error{
			Type:    types.TypeErrorUrlNotFound,
			Code:    types.CodeErrorUrlNotFound,
			Message: log,
		}
		return nil, c.Error
	}

	filename, ext, contentType := "", "", ""

	contentType = resp.Header.Get("Content-Type")

	// get file extension
	if filename == "" {
		filename = "tmp"
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil && len(exts) > 0 {
			ext = strings.ToLower(exts[len(exts)-1])
		}
	} else {
		ext = path.Ext(filename)
	}

	nameFile := filename + ext

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	if err := writer.WriteField("messaging_product", msg.GetMessagingProduct()); err != nil {
		c.Error = fmt.Errorf("Error in multiparRequest when write messaging_product. Error is: %s", err.Error())
		return nil, c.Error
	}

	if err := writer.WriteField("type", contentType); err != nil {
		c.Error = fmt.Errorf("Error in multiparRequest when write type. Error is: %s", err.Error())
		return nil, c.Error
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", nameFile))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, nil
	}

	if _, err = io.Copy(part, resp.Body); err != nil {
		c.Error = fmt.Errorf("Error in multiparRequest when copy file to form. Error is: %s", err.Error())
		return nil, c.Error
	}

	urlPath, _ = url.Parse(c.path + ePoint)
	urlPath = c.BaseUrl.ResolveReference(urlPath)
	c.request, e = http.NewRequest(methoth, urlPath.String(), payload)
	if e != nil {
		c.Error = fmt.Errorf("Error in multipartRequest, NewRequest: %s. Error is: %s", c.BaseUrl, e.Error())
		return nil, c.Error
	}

	multiparHeader(c, writer.FormDataContentType())
	return c.request, nil
}
