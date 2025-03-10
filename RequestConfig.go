//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ecsavigne/client_wa_oficial/types"
)

func deafaultHeader(c *Config) {
	c.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	c.request.Header.Set("Content-Type", "application/json")
}

func multiparHeader(c *Config, m *multipart.Writer) {
	c.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	c.request.Header.Set("Content-Type", m.FormDataContentType())
}

func deafaultRequest(methoth string, ePoint string, c *Config, msg types.Messager) (*http.Request, error) {
	var e error
	var urlPath *url.URL

	if !strings.HasPrefix(ePoint, "/") {
		var log = "Error in deafultRequest, file: RequestConfig.go.Error is: EndPoint is not start with /"
		fmt.Println(log)
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
		var log = "Error in deafultRequest, file: RequestConfig.go.Error is: EndPoint is not start with /"
		fmt.Println(log)
		c.Error = fmt.Errorf("%s", log)
		panic(c.Error)
	}

	urlPath, _ = url.Parse(c.path + ePoint)
	urlPath = c.BaseUrl.ResolveReference(urlPath)

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer pw.Close()

		// Obtener informacion del archivo
		resp, err := http.Get(msg.GetMessageLink())
		if err != nil {
			pw.CloseWithError(err)
			c.Error = fmt.Errorf("Error in multiparRequest getting file. Error is: %s", e.Error())
			return
		}
		defer resp.Body.Close()

		// get name
		filename, ext := "", ""
		if disp := resp.Header.Get("Content-Disposition"); disp != "" {
			_, params, err := mime.ParseMediaType(disp)
			if err == nil {
				filename = params["filename"]
			}
		}

		// get file extension
		if filename == "" {
			filename = "tmp"
			// Obtener typeMime del archivo
			contentType := resp.Header.Get("Content-Type")
			exts, err := mime.ExtensionsByType(contentType)
			if err == nil && len(exts) > 0 {
				ext = strings.ToLower(exts[0])
			}
		} else {
			ext = path.Ext(filename)
		}

		// Crea o abre el archivo donde se almacenará el contenido
		nameFile := filename + "." + ext
		// file, err := os.Create(nameFile)
		// if err != nil {
		// 	log.Fatalf("Error in creating file: %v", err)
		// }
		// defer file.Close()

		// // Copy to path local
		// if _, err = io.Copy(file, resp.Body); err != nil {
		// 	pw.CloseWithError(err)
		// 	return
		// }

		if err := writer.WriteField("messaging_product", msg.GetMessagingProduct()); err != nil {
			pw.CloseWithError(err)
			c.Error = fmt.Errorf("Error in multiparRequest when write messaging_product. Error is: %s", e.Error())
			return
		}

		// Crear el campo del archivo en el formulario
		part, err := writer.CreateFormFile("file", nameFile)
		if err != nil {
			pw.CloseWithError(err)
			c.Error = fmt.Errorf("Error in multiparRequest when create form file. Error is: %s", e.Error())
			return
		}

		// Copiar el contenido del archivo al formulario
		// if _, err = io.Copy(part, file); err != nil {
		if _, err = io.Copy(part, resp.Body); err != nil {
			pw.CloseWithError(err)
			c.Error = fmt.Errorf("Error in multiparRequest when copy file to form. Error is: %s", e.Error())
			return
		}

		// close body multipart
		if err = writer.Close(); err != nil {
			pw.CloseWithError(err)
			c.Error = fmt.Errorf("Error in multiparRequest when close body multipart. Error is: %s", e.Error())
			return
		}
	}()

	c.request, e = http.NewRequest(methoth, urlPath.String(), pr)
	if e != nil {
		c.Error = fmt.Errorf("Error in multipartRequest, NewRequest: %s. Error is: %s", c.BaseUrl, e.Error())
		return nil, c.Error
	}

	multiparHeader(c, writer)
	return c.request, nil
}
