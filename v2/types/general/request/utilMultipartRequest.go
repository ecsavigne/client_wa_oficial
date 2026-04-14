package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path"
	"strings"

	clientoficial "github.com/ecsavigne/client_wa_oficial/v2/client"
	generalpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	"google.golang.org/protobuf/proto"
)

func UtilMultipartRequest(cl clientoficial.Client, msg proto.Message, method, ePoint string) (*http.Request, error) {
	var (
		e                          error
		urlPath                    *url.URL
		request                    *http.Request
		filename, ext, contentType = "", "", ""
		fileTemp                   *bytes.Buffer
		// resp                       *http.Response
		// client                     = createClientHttp2(60)
		fileDescriptor, ok = msg.(*generalpbv1.FileDescriptor)
		data               = struct {
			ePoint, method string
		}{ePoint: ePoint, method: method}
	)

	if ok && fileDescriptor != nil {
		filename = fileDescriptor.GetFilename()
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil && len(exts) > 0 {
			ext = strings.ToLower(exts[len(exts)-1])
		}
		if !strings.HasSuffix(filename, "."+ext) {
			filename += "." + ext
		}
		contentType = fileDescriptor.GetMimeType()
		fileTemp = bytes.NewBuffer(fileDescriptor.GetContent())
	} else {
		return nil, fmt.Errorf("Error in multipartRequest, file: RequestConfig.go.Error is: FileDescriptor is nil")
		// resp, e = client.Get(data.link)
		// if e != nil {
		// 	e = fmt.Errorf("Error in multiparRequest getting file. Error is: %s", e.Error())
		// 	return nil, e
		// }
		// defer resp.Body.Close()

		// contentType = resp.Header.Get("Content-Type")

		// switch resp.StatusCode {
		// case 400:
		// 	log := fmt.Sprintf("Error in function makeRequest bad request of %s. Message type: %s. status: %s", data.typeClient, data.msgType, resp.Status)
		// 	e = fmt.Errorf("type: %s, code: %d, fileInfo: %s", types.MsgErrorBadRequest, log)
		// 	return nil, e
		// case 401:
		// 	log := fmt.Sprintf("Error in function makeRequest bad request of %s. Message type: %s. s: %s", data.typeClient, data.msgType, resp.Status)
		// 	e = fmt.Errorf("type: %s, code: %d, fileInfo: %s", types.MsgErrorUnauthorized, log)
		// 	return nil, e
		// case 404:
		// 	log := fmt.Sprintf("Error in function makeRequest bad request of %s. Message type: %s. s: %s", data.typeClient, data.msgType, resp.Status)
		// 	e = fmt.Errorf("type: %s, code: %d, fileInfo: %s", types.MsgErrorUrlNotFound, log)
		// 	return nil, e
		// }
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	switch cl.GetType() {
	case clientoficial.CLIENT_WHATSAPP.String():
		if err := writer.WriteField("messaging_product", "whatsapp"); err != nil {
			e = fmt.Errorf("Error in multiparRequest when write messaging_product. Error is: %s", err.Error())
			return nil, e
		}
	case clientoficial.CLIENT_IG.String():
		// if err := writer.WriteField("messaging_product", data.fileInfo.GetMessagingProduct()); err != nil {
		// 	e = fmt.Errorf("Error in multiparRequest when write messaging_product. Error is: %s", err.Error())
		// 	return nil, e
		// }
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("Error in multiparRequest when create part. Error is: %s", err.Error())
	}

	// copy bin
	if _, err = io.Copy(part, fileTemp); err != nil {
		e = fmt.Errorf("Error in multiparRequest when copy file to form. Error is: %s", err.Error())
		return nil, e
	}

	if e := writer.Close(); e != nil {
		e = fmt.Errorf("Error closing multipart writer: %s", e.Error())
		return nil, e
	}

	urlPath, _ = url.Parse(data.ePoint)
	baseUrl, _ := url.Parse(path.Join(cl.GetConfig().GetBaseUrl(), cl.GetConfig().GetVersion()))
	urlPath = baseUrl.ResolveReference(urlPath)

	request, e = http.NewRequestWithContext(context.Background(), data.method, urlPath.String(), payload)
	if e != nil {
		e = fmt.Errorf("Error in multipartRequest, NewRequest: %s. Error is: %s", baseUrl, e.Error())
		return nil, e
	}

	multiparHeader(cl.GetConfig(), request, writer.FormDataContentType())

	return request, nil
}
