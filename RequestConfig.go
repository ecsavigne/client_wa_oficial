//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path"
	"strings"

	"github.com/ecsavigne/client_wa_oficial/types"
	"golang.org/x/net/http2"
)

func deafaultHeader(c *Config) {
	c.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	c.request.Header.Set("Content-Type", "application/json")
}

func multiparHeader(c *Config, contentType string) {
	c.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	c.request.Header.Set("Content-Type", contentType)
}

func deafaultRequest(methoth string, ePoint string, c *Config, msg types.Messager) (*http.Request, error) {
	var e error
	var urlPath *url.URL

	if !strings.HasPrefix(ePoint, "/") {
		var log = "Error in deafultRequest, file: RequestConfig.go.Error is: EndPoint is not start with /"
		c.Error = fmt.Errorf("%s", log)
		panic(c.Error)
	}

	urlPath, _ = url.Parse(c.path + ePoint)
	urlPath = c.BaseUrl.ResolveReference(urlPath)
	c.request, e = http.NewRequest(methoth, urlPath.String(), msg.ToJSONReader())
	if e != nil {
		c.Error = fmt.Errorf("Error in deafaultRequest, NewRequest: %s. Error is: %s", c.BaseUrl, e.Error())
		return nil, c.Error
	}
	deafaultHeader(c)
	return c.request, nil
}

func multipartRequest(methoth string, ePoint string, c *Config, msg types.Messager) (*http.Request, error) {
	var e error
	var urlPath *url.URL
	if !strings.HasPrefix(ePoint, "/") {
		var log = "Error in multipartRequest, file: RequestConfig.go.Error is: EndPoint is not start with /"
		c.Error = fmt.Errorf("%s", log)
		panic(c.Error)
	}

	tr := &http.Transport{}
	err := http2.ConfigureTransport(tr)
	if err != nil {
		c.Error = fmt.Errorf("Error configuring HTTP/2. Error is: %s", err.Error())
		return nil, c.Error
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get(msg.GetMessageLink())
	if err != nil {
		c.Error = fmt.Errorf("Error in multiparRequest getting file. Error is: %s", err.Error())
		return nil, c.Error
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		c.Error = fmt.Errorf("Error code is 404. Error is: %s", resp.Status)
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
