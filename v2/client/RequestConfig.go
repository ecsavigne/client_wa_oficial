//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"bytes"
	"context"
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
	"time"

	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/message"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
	"golang.org/x/net/http2"
)

type TypeRequest = string

const (
	RequestGetMessageInfo    TypeRequest = "RequestGetMessageInfo"
	RequestDeleteMedia       TypeRequest = "RequestDeleteMedia"
	RequestChangeUrlFull     TypeRequest = "RequestChangeUrlFull"
	RequestWithQueryPhone    TypeRequest = "RequestWithQueryPhone"
	RequestWithQuery         TypeRequest = "RequestWithQuery"
	RequestWithQueryBusiness TypeRequest = "RequestWithQueryBusiness"
	RequestWithVersion       TypeRequest = "RequestWithVersion"
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

func (q QueryData) SetValue(key string, value any) {
	q[key] = value
}

func (q QueryData) GetValue(key string) any {
	return q[key]
}

func defaultHeader(c *Config, contentType ...string) {
	cT := "application/json"
	if len(contentType) > 0 {
		cT = contentType[0]
	}

	h := make(http.Header)
	h.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	h.Set("Content-Type", cT)
	c.request.Header = h
}

func multiparHeader(c *Config, contentType string) {
	h := make(http.Header)
	h.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	h.Set("Content-Type", contentType)
	c.request.Header = h
}

func resetError(c *Config) {
	if c.Error != nil {
		c.Error = nil
	}
}

// defaultRequest creates a new request with the given method, endpoint, and optional parameters.
// If the endpoint does not start with a slash, it will be appended to the configured base URL.
// If the first parameter is a message.Messager, it will be marshaled to JSON and used as the request body.
// If the first parameter is a map[string]any, it will be marshaled to JSON and used as the request body.
// If the first parameter is a TypeRequest, it will be used to determine how to construct the request URL.
// The following types are supported:
//   - RequestGetMessageInfo: The endpoint will be appended to the configured versioned path.
//   - RequestDeleteMedia: The endpoint will be appended to the configured path.
//   - RequestWithVersion: The endpoint will be appended to the configured versioned path.
//   - RequestWithQuery: The endpoint will be appended to the configured path, and the second parameter
//     (which must be a QueryData) will be used to construct the query string.
//   - RequestWithQueryBusiness: The endpoint will be appended to the configured business path, and the second parameter
//     (which must be a QueryData) will be used to construct the query string.
//
// If the first parameter is not one of the above types, an error will be returned.
// The second parameter (if present) is used to construct the query string if the first parameter is a TypeRequest.
// The third parameter (if present) is used to construct the request body if the first parameter is a message.Messager.
// The function returns the created request, a cancel function, and an error (if any).
func defaultRequest(methoth string, ePoint string, c *Config, params ...any) (*http.Request, context.CancelFunc, error) {
	resetError(c)
	var (
		e              error
		urlPath        *url.URL
		urlAlternative string = ""
		msg            message.Messager
		formData       *bytes.Buffer = bytes.NewBuffer([]byte{})
		ctx            context.Context
		cancel         context.CancelFunc
	)

	if !strings.HasPrefix(ePoint, "/") {
		var log = "Error in deafultRequest, file: RequestConfig.go.Error is: EndPoint is not start with /"
		c.Error = fmt.Errorf("%s", log)
		return nil, nil, c.Error
	}

	urlPath, _ = url.Parse(c.path + ePoint)
	urlPath = c.BaseUrl.ResolveReference(urlPath)

	if len(params) > 0 && len(params) == 1 {
		switch v := params[0].(type) {
		case message.Messager:
			msg = v
			c.request, e = http.NewRequest(methoth, urlPath.String(), msg.ToJSONReader())
		case map[string]any:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, nil, response.NewError(&response.Error{
					Type:    types.TypeErrorUnrecognized,
					Code:    types.CodeErrorUnrecognized,
					Message: err.Error(),
				})
			}
			formData = bytes.NewBuffer(b)
			c.request, e = http.NewRequest(methoth, urlPath.String(), formData)
		case TypeRequest:
			switch v {
			case RequestGetMessageInfo, RequestDeleteMedia, RequestWithVersion:
				urlPath, _ = url.Parse(fmt.Sprintf("%s%s", c.pathVersion, ePoint))
				urlPath = c.BaseUrl.ResolveReference(urlPath)
				c.request, e = http.NewRequest(methoth, urlPath.String(), nil)
			case RequestChangeUrlFull:
				urlAlternative = strings.TrimPrefix(ePoint, "/")
				ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
				c.request, e = http.NewRequestWithContext(ctx, methoth, urlAlternative, nil)
			}
		}
	} else {
		if len(params) > 2 {
			c.Error = response.NewError(&response.Error{
				Type:    types.TypeErrorUnrecognized,
				Code:    types.CodeErrorUnrecognized,
				Message: "Error in deafultRequest, file: RequestConfig.go. Context: len(params) > 2",
			})
			return nil, nil, c.Error
		}

		queryData := ""
		if obj, ok := params[1].(QueryData); ok {
			queryData = obj.String()
		}

		switch v := params[0].(type) {
		case TypeRequest:
			switch v {
			case RequestWithQueryBusiness:
				urlPath, _ = url.Parse(fmt.Sprintf("%s%s?%s", c.pathBusiness, ePoint, queryData))
			case RequestWithQueryPhone:
				urlPath, _ = url.Parse(fmt.Sprintf("%s%s?%s", c.path, ePoint, queryData))
			case RequestWithQuery:
				urlPath, _ = url.Parse(fmt.Sprintf("%s%s?%s", c.pathVersion, ePoint, queryData))
			case RequestWithVersion:
				b, err := json.Marshal(params[1])
				if err != nil {
					return nil, nil, response.NewError(&response.Error{
						Type:    types.TypeErrorUnrecognized,
						Code:    types.CodeErrorUnrecognized,
						Message: err.Error(),
					})
				}
				formData = bytes.NewBuffer(b)
				urlPath, _ = url.Parse(fmt.Sprintf("%s%s", c.pathVersion, ePoint))
			}
		}

		urlPath = c.BaseUrl.ResolveReference(urlPath)
		c.request, e = http.NewRequest(methoth, urlPath.String(), formData)
	}

	if e != nil {
		if cancel != nil {
			cancel()
		}
		c.Error = fmt.Errorf("Error in defaultRequest, NewRequest: %s. Error is: %s", c.BaseUrl, e.Error())
		return nil, nil, c.Error
	}
	defaultHeader(c)
	return c.request, cancel, nil
}

func multipartRequest(methoth string, ePoint string, c *Config, msg message.Messager) (*http.Request, error) {
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

	var resp *http.Response
	filename, ext, contentType := "", "", ""
	var fileTemp multipart.File

	if msg.GetFileHeader() != nil {
		filename = msg.GetFileHeader().Filename
		contentType = msg.GetFileHeader().Header.Get("Content-Type")
		ext = path.Ext(filename)
		fileTemp, e = msg.GetFileHeader().Open()
		if e != nil {
			c.Error = fmt.Errorf("Error in multiparRequest getting file. Error is: %s", e.Error())
			return nil, c.Error
		}
		defer fileTemp.Close()
	} else {
		resp, e = client.Get(msg.GetMessageLink())
		if e != nil {
			c.Error = fmt.Errorf("Error in multiparRequest getting file. Error is: %s", e.Error())
			return nil, c.Error
		}
		defer resp.Body.Close()

		contentType = resp.Header.Get("Content-Type")

		switch resp.StatusCode {
		case 400:
			log := fmt.Sprintf("Error in function makeRequest bad request of ClientWA. Message type: %s. error is: %s", msg.GetType(), resp.Status)
			c.Error = response.NewError(&response.Error{
				Type:    types.TypeErrorBadRequest,
				Code:    types.CodeErrorBadRequest,
				Message: log,
			})
			return nil, c.Error
		case 401:
			log := fmt.Sprintf("Error in function makeRequest bad request of ClientWA. Message type: %s. error is: %s", msg.GetType(), resp.Status)
			c.Error = response.NewError(&response.Error{
				Type:    types.TypeErrorUnauthorized,
				Code:    types.CodeErrorUnauthorized,
				Message: log,
			})
			return nil, c.Error
		case 404:
			log := fmt.Sprintf("Error in function makeRequest bad request of ClientWA. Message type: %s. error is: %s", msg.GetType(), resp.Status)
			c.Error = response.NewError(&response.Error{
				Type:    types.TypeErrorUrlNotFound,
				Code:    types.CodeErrorUrlNotFound,
				Message: log,
			})
			return nil, c.Error
		}
	}

	// get file extension
	if filename == "" {
		filename = "tmp"
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil && len(exts) > 0 {
			ext = strings.ToLower(exts[len(exts)-1])
		}
		filename += ext
	}

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
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, nil
	}

	if msg.GetFileHeader() != nil {
		if _, err = io.Copy(part, fileTemp); err != nil {
			c.Error = fmt.Errorf("Error in multiparRequest when copy file to form. Error is: %s", err.Error())
			return nil, c.Error
		}
		msg.ResetFileHeader()
	} else {
		if _, err = io.Copy(part, resp.Body); err != nil {
			c.Error = fmt.Errorf("Error in multiparRequest when copy file to form. Error is: %s", err.Error())
			return nil, c.Error
		}
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
