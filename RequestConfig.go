//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ecsavigne/client_wa_oficial/types"
)

func deafaultHeader(c *Config) {
	c.request.Header.Set("Content-Type", "application/json")
	c.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
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
