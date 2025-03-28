//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ecsavigne/client_wa_oficial/v2/event"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"golang.org/x/net/http2"

	"github.com/gorilla/websocket"
)

type clientHttp struct {
	*http.Client
	BaseUrl *url.URL `json:"base_url"`
}

type Config struct {
	Token               string `json:"token"`
	WaBusinessAccountId string `json:"wa_business_account_id"`
	WaPhoneNumberId     string `json:"wa_phone_number_id"`
	// Path del archivo .env incluyendo el nombre del archivo sin la extensión ej: file: /.../../config_env.env -> EnvFilePath: /.../../config_env
	EnvFilePath string `json:"env_file_path"`
	Error       error
	// Url del servidor WebHook con ruta /ws para conectar con el servidor WebSocket ej: wss://webhooks.savcoe-services.com/ws
	WebhookSocket string    `json:"webhook_socket"`
	EventHandle   func(any) // Funcion para manejar los eventos del servidor WebHook WebSocket
	path          string
	pathVersion   string
	pathBusiness  string
	clientHttp
	request   *http.Request
	MediaInfo *types.MediaInfo
}

type ClientWA struct {
	*Config `json:"config"`
}

func (cl *ClientWA) initWebHookSocket() {
	// url = "wss://webhooks.savcoe-services.com/ws"
	// Conectar al servidor WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(cl.Config.WebhookSocket, nil)
	if err != nil {
		evt := &event.EventErrorSocketConnect{}
		switch {
		case websocket.IsUnexpectedCloseError(err):
			evt.Error = types.Error{
				Type:    types.TypeErrorUnexpectedClose,
				Code:    types.CodeErrorUnexpectedClose,
				Message: types.MsgErrorUnexpectedClose,
			}
		case strings.Contains(err.Error(), "tls: internal error"):
			evt.Error = types.Error{
				Type:    types.TypeErrorTlsInternal,
				Code:    types.CodeErrorTlsInternal,
				Message: types.MsgErrorTlsInternal,
			}
		case strings.Contains(err.Error(), "bad handshake"):
			evt.Error = types.Error{
				Type:    types.TypeErrorBadHandshake,
				Code:    types.CodeErrorBadHandshake,
				Message: types.MsgErrorBadHandshake,
			}
		case strings.Contains(err.Error(), "dial tcp: lookup ws"):
			evt.Error = types.Error{
				Type:    types.TypeErrorDialTcp,
				Code:    types.CodeErrorDialTcp,
				Message: types.MsgErrorDialTcp,
			}
		default:
			evt.Error = types.Error{
				Type:    types.TypeErrorUnrecognizedWebSocket,
				Code:    types.CodeErrorUnrecognizedWebSocket,
				Message: fmt.Sprintf("%s. Original error: %s", types.MsgErrorUnrecognizedWebSocket, err.Error()),
			}
		}
		cl.EventHandle(evt)
		return
	}
	defer conn.Close()

	// Listener message of the server way socket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Conexión cerrada:", err)
			break
		}
		cl.Config.EventHandle(message)
	}
}

// Create one Client of WhatsApp Official return *ClientWA :
// - If the EnvFilePath or Path in Config.EnvFilePath not found. ClientWA.Error = &types.Error{Type: types.TypeErrorConfig, Code: types.CodeErrorEnvNotFound, Message: types.MsgErrorEnvNotFound}
// - If occurred error in conection with WebHook Socket emit one event type: event.EventErrorSocketConnect
func NewClientWA(c ...Config) *ClientWA {
	if len(c) == 0 {
		c = append(c, Config{})
	}

	cl := &ClientWA{
		Config: &c[0],
	}
	err := setEnv(c[0])
	if err != nil {
		cl.Error = err
		return cl
	}

	cl.Config = newConfig(c[0])
	if cl.Error != nil {
		return cl
	}

	if c[0].WebhookSocket != "" && c[0].EventHandle != nil {
		go cl.initWebHookSocket()
	}

	return cl
}

func newConfig(c Config) *Config {
	c.Error = nil
	if WA_BASE_URL == "" {
		c.Error = &types.Error{
			Type:    types.TypeErrorBaseUrlEmpty,
			Code:    types.CodeErrorBadHandshake,
			Message: types.MsgErrorBaseUrlEmpty,
		}
		return &c
	}

	if CLOUD_API_VERSION == "" {
		c.Error = &types.Error{
			Type:    types.TypeErrorApiVersionEmpty,
			Code:    types.CodeErrorApiVersionEmpty,
			Message: types.MsgErrorApiVersionEmpty,
		}
		return &c
	}

	if WA_PHONE_NUMBER_ID == "" {
		c.Error = &types.Error{
			Type:    types.TypeErrorPhoneIdEmpty,
			Code:    types.CodeErrorPhoneIdEmpty,
			Message: types.MsgErrorPhoneIdEmpty,
		}
		return &c
	}

	if WA_BUSINESS_ACCOUNT_ID == "" {
		c.Error = &types.Error{
			Type:    types.TypeErrorBusinessIdEmpty,
			Code:    types.CodeErrorBusinessIdEmpty,
			Message: types.MsgErrorBusinessIdEmpty,
		}
		return &c
	}

	c.path = path.Join(CLOUD_API_VERSION, WA_PHONE_NUMBER_ID)
	c.pathBusiness = path.Join(CLOUD_API_VERSION, WA_BUSINESS_ACCOUNT_ID)
	c.pathVersion = path.Join(CLOUD_API_VERSION)

	c.BaseUrl, _ = url.Parse(WA_BASE_URL)

	if c.Token == "" {
		// c.Token = CLOUD_API_ACCESS_TOKEN
		c.Error = &types.Error{
			Type:    types.TypeErrorTokenEmpty,
			Code:    types.CodeErrorTokenEmpty,
			Message: types.MsgErrorTokenEmpty,
		}
		return &c
	}

	if c.Client == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12},
		}
		http2.ConfigureTransport(tr)

		c.Client = &http.Client{
			Timeout:   30 * time.Second,
			Transport: tr,
		}
	}

	return &c
}

func (c *ClientWA) resetMessageInfo() {
	if c.MediaInfo != nil {
		c.MediaInfo = nil
	}
}

func doRequest(req *http.Request, c *ClientWA) (*http.Response, error) {
	res, e := c.clientHttp.Do(req)
	if e != nil {
		log := fmt.Sprintf("Error in function doRequest when send HTTP request to server with Do. Error is: %s", e.Error())
		c.Config.Error = fmt.Errorf("%s", log)
		return nil, c.Config.Error
	}

	switch res.StatusCode {
	case 400:
		log := fmt.Sprintf("Error in function doRequest. Code: %d, Message: %s, MetaError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"))
		c.Config.Error = &types.Error{
			Type:    types.TypeErrorBadRequest,
			Code:    types.CodeErrorBadRequest,
			Message: log,
		}
		return nil, c.Config.Error
	case 401:
		log := fmt.Sprintf("Error in function doRequest. Code: %d, Message: %s, MetaError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"))
		c.Config.Error = &types.Error{
			Type:    types.TypeErrorUnauthorized,
			Code:    types.CodeErrorUnauthorized,
			Message: log,
		}
		return nil, c.Config.Error
	case 404:
		log := fmt.Sprintf("Error in function doRequest. Code: %d, Message: %s, MetaError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"))
		c.Config.Error = &types.Error{
			Type:    types.TypeErrorUrlNotFound,
			Code:    types.CodeErrorUrlNotFound,
			Message: log,
		}
		return nil, c.Config.Error
	}

	return res, nil
}

func (c *ClientWA) doRequest(req *http.Request) (types.ResponserRequest, error) {
	res, e := doRequest(req, c)

	var responser types.ResponserRequest
	if e != nil {
		responser = &types.Error{
			Type:    types.TypeErrorInRequest,
			Code:    types.CodeErrorInRequest,
			Message: fmt.Sprintf("Type: %s. Error is: %s", types.MsgErrorInRequest, e.Error()),
		}
		return responser, c.Config.Error
	}

	bodyResponse, e := io.ReadAll(res.Body)
	if e != nil {
		log := fmt.Sprintf("Error in function doRequest of ClientWA when reading response body. Error is: %s", e.Error())
		responser = &types.Error{
			Type:    types.TypeErrorUnrecognized,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Type: %s. Error is: %s", types.MsgErrorUnrecognized, e.Error()),
		}
		c.Config.Error = fmt.Errorf("%s", log)
		return nil, c.Config.Error
	}

	defer res.Body.Close()

	return types.JsonWrapperResponseRequest(bodyResponse), nil
}

func (c *ClientWA) makeRequest(methoth string, ePoint string, msg types.Messager) (types.ResponserRequest, error) {
	if msg.GetMessageLink() != "" {
		multipartRequest(methoth, ePoint, c.Config, msg)
	} else {
		defaultRequest(methoth, ePoint, c.Config, msg)
	}

	if c.Config.Error != nil {
		log := fmt.Sprintf("Error in function makeRequest creting request of ClientWA. error is: %s", c.Config.Error.Error())
		c.Config.Error = fmt.Errorf("%s", log)
		return nil, c.Config.Error
	}

	var (
		responseReq types.ResponserRequest
		e           error
	)

	responseReq, e = c.doRequest(c.request)
	if responseReq.GetResponseError() != nil {
		return responseReq.GetResponseError(), nil
	}

	if e != nil {
		return nil, c.Config.Error
	}

	return responseReq, nil
}

func validTypeMsg(msg types.Messager, msgType string) bool {
	return msg.GetType() == msgType
}

func validTypeInteractive(msg types.MessageInteractive, interactiveType string) bool {
	return msg.InteractiveProto.Type == interactiveType
}

// SendTemplate sends a template message. It validates the message type to ensure it is a template.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendTemplate(m types.Messager) types.ResponserRequest {

	if !validTypeMsg(m, types.MessageTypeTemplate) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeTemplate, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendTemplate request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendTextMessage sends a text message. It validates the message type to ensure it is a text.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendTextMessage(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeText) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeText, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendText request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendReaction sends a reaction message. It validates the message type to ensure it is a reaction.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendReaction(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeReaction) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeReaction, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendReaction request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendInteractiveList sends an interactive list message. It validates the message type to ensure it is an interactive of type list.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendInteractiveList(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m.(types.MessageInteractive), types.MessageTypeInteractive) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		}
	} else if !validTypeInteractive(m.(types.MessageInteractive), types.InteractiveTypeList) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeList),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveList request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendInteractiveButtonResponse sends an interactive button response message. It validates the message type to ensure it is an
// interactive of type button response. If the message type is incorrect, it returns an error. Otherwise, it makes a request to
// send the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendInteractiveButtonResponse(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m.(types.MessageInteractive), types.MessageTypeInteractive) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		}
	} else if !validTypeInteractive(m.(types.MessageInteractive), types.InteractiveTypeButtonResponse) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeButtonResponse),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveButtonResponse request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendInteractiveButtonUrl sends an interactive button URL message. It validates the message type to ensure it is an
// interactive of type button URL. If the message type is incorrect, it returns an error. Otherwise, it makes a request
// to send the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendInteractiveButtonUrl(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m.(types.MessageInteractive), types.MessageTypeInteractive) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		}
	} else if !validTypeInteractive(m.(types.MessageInteractive), types.InteractiveTypeButtonUrl) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeButtonUrl),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveButtonUrl request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendInteractiveMsgProcess sends an interactive process message. It validates the message type to ensure it is an
// interactive of type process. If the message type is incorrect, it returns an error. Otherwise, it makes a request
// to send the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendInteractiveMsgProcess(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m.(types.MessageInteractive), types.MessageTypeInteractive) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		}
	} else if !validTypeInteractive(m.(types.MessageInteractive), types.InteractiveTypeProcess) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeProcess),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveMsgProcess request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendInteractiveOneProduct sends an interactive message of type product. It validates the message type to ensure it is an
// interactive of type product. If the message type is incorrect, it returns an error. Otherwise, it makes a request
// to send the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendInteractiveOneProduct(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m.(types.MessageInteractive), types.MessageTypeInteractive) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		}
	} else if !validTypeInteractive(m.(types.MessageInteractive), types.InteractiveTypeProduct) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeProduct),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveOneProduct request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendInteractiveMultiProduct sends an interactive message of type multi product. It validates the message type to
// ensure it is an interactive of type multi product. If the message type is incorrect, it returns an error. Otherwise,
// it makes a request to send the message. If the request is successful, the response is returned; otherwise, an error
// is returned.
func (c *ClientWA) SendInteractiveMultiProduct(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m.(types.MessageInteractive), types.MessageTypeInteractive) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		}
	} else if !validTypeInteractive(m.(types.MessageInteractive), types.InteractiveTypeMultiProduct) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeMultiProduct),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveMultiProduct request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendInteractiveCatalog sends an interactive catalog message. It validates the message type
// to ensure it is an interactive of type catalog. If the message type is incorrect, it returns
// an error. Otherwise, it makes a request to send the message. If the request is successful,
// the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendInteractiveCatalog(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m.(types.MessageInteractive), types.MessageTypeInteractive) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		}
	} else if !validTypeInteractive(m.(types.MessageInteractive), types.InteractiveTypeCatalog) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeCatalog),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveCatalog request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendResponseMsg sends a response message. It validates the message type to ensure it is a response
// message. If the message type is incorrect, it returns an error. Otherwise, it makes a request to send
// the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendResponseMsg(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, "response") {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", "response", m.GetType()),
		}
	} else {
		if m.IsTypeResponse() {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message type response not expect, type: '%s'", m.(*types.MessageResponse).Header.Type),
			}
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendResponseMsg request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

func (c *ClientWA) validLinAndId(m types.Messager) types.ResponserRequest {
	if m.GetMessageLink() != "" && m.GetMessageId() != "" {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: "Expect Message.id or Message.link, but not both",
		}
	}

	return nil
}

// SendAudioMessage sends an audio message. It validates the message type to ensure it is an audio.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendAudioMessage(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeAudio) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeAudio, m.GetType()),
		}
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" {
		r := c.UploadFile(m, types.AUDIO)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendAudio request of ClientWA. error is: ", e.Error()),
		}
	}

	if resp.GetType() == types.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendImageMessage sends an image message. It validates the message type to ensure it is an image.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendImageMessage(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeImage) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeImage, m.GetType()),
		}
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" {
		r := c.UploadFile(m, types.IMAGE)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendImage request of ClientWA. error is: ", e.Error()),
		}
	}

	if resp.GetType() == types.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendVideoMessage sends a video message. It validates the message type to ensure it is a video.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendVideoMessage(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeVideo) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeVideo, m.GetType()),
		}
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" {
		r := c.UploadFile(m, types.VIDEO)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendVideo request of ClientWA. error is: ", e.Error()),
		}
	}

	if resp.GetType() == types.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendDocumentMessage sends a document message. It validates the message type to ensure it is a document.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendDocumentMessage(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeDocument) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeDocument, m.GetType()),
		}
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" {
		r := c.UploadFile(m, types.DOCUMENT)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendDocument request of ClientWA. error is: ", e.Error()),
		}
	}

	if resp.GetType() == types.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendStickerMessage sends a sticker message. It validates the message type to ensure it is a sticker. If the message type is
// incorrect, it returns an error. Otherwise, it makes a request to send the message. If the request is successful, the response
// is returned; otherwise, an error is returned.
func (c *ClientWA) SendStickerMessage(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeSticker) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeSticker, m.GetType()),
		}
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" {
		r := c.UploadFile(m, types.STICKER)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendSticker request of ClientWA. error is: ", e.Error()),
		}
	}

	if resp.GetType() == types.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendLocationMessage sends a location message. It validates the message type
// to ensure it is a location message. If the message type is incorrect, it
// returns an error. Otherwise, it makes a request to send the message. If the
// request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendLocationMessage(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeLocation) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeLocation, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendLocationMessage request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// SendContactMessage sends a contact message. It first validates the message type to ensure it is a contact message.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) SendContactMessage(m types.Messager) types.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeContact) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeContact, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*types.Error); ok != nil {
			return e.(*types.Error)
		}

		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendContactMessage request of ClientWA. error is: ", e.Error()),
		}
	}

	return resp
}

// UploadFile uploads a file to the server. It validates the message type to ensure it is either a text, audio, image, video,
// document, or sticker message. If the message type is incorrect, it returns an error. Otherwise, it makes a request to upload
// the file. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) UploadFile(m types.Messager, mt types.MediaType) types.ResponserRequest {
	resp, e := c.makeRequest(http.MethodPost, "/media", m)
	if e != nil {
		msgError := fmt.Sprintln("Error in UploadFile request of ClientWA. error is: ", e.Error())
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: msgError,
		}

	} else if resp.GetType() == types.ResponseError {
		return resp.GetResponseError()
	}

	id := ""

	if resp.GetType() == types.ResponseMediaInfo {
		c.Config.MediaInfo = resp.GetResponseMediaInfo()
		id = c.Config.MediaInfo.ID
	}

	switch mt {
	case types.AUDIO, types.IMAGE, types.VIDEO, types.DOCUMENT, types.STICKER:
		m.SetId(id)
		m.SetLink("")
	}

	return resp
}

// DownloadFile downloads a file using its unique identifier. It first retrieves the file information
// and checks if the response indicates an error or contains media information. If media information
// is available, it sends a request to download the file and saves it to the specified path with the
// given name. If an error occurs during any step, it returns an error. Otherwise, it completes the
// download process successfully.
func (c *ClientWA) DownloadFile(id, path, nameFile string) error {
	responseReq, e := c.getFileInfo(id)
	if e != nil {
		return c.Config.Error
	}

	if responseReq.IsType(types.ResponseError) {
		return responseReq.GetResponseError()
	}

	if responseReq.IsType(types.ResponseMediaInfo) {
		// Get binaryFile
		mInfo := responseReq.GetResponseMediaInfo()
		_, cancel, _ := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", mInfo.Url), c.Config, RequestChangeUrlFull)

		if c.Config.Error != nil {
			return c.Config.Error
		} else {
			// Save binaryFile in path
			res, e := doRequest(c.request, c)
			ext := strings.Split(res.Header.Get("Content-Disposition"), ".")
			if e != nil {
				return e
			} else {
				file, e := os.Create(fmt.Sprintf("%s%s.%s", path, nameFile, ext[len(ext)-1]))
				if e != nil {
					return e
				}
				defer file.Close()

				_, e = io.Copy(file, res.Body)
				if e != nil {
					return &types.Error{
						Type:    types.ResponseError,
						Code:    types.CodeErrorUnrecognized,
						Message: fmt.Sprintln("Error in DownloadFile request of ClientWA. error is: ", e.Error()),
					}
				}
				defer res.Body.Close()
				defer cancel()
			}
		}
	}
	return nil
}

func (c *ClientWA) getFileInfo(id string) (types.ResponserRequest, error) {
	// Crear request
	defaultRequest(http.MethodGet, fmt.Sprintf("/%s", id), c.Config, RequestGetMessageInfo)

	var (
		responseReq types.ResponserRequest
		e           error
	)

	responseReq, e = c.doRequest(c.request)
	if responseReq.GetResponseError() != nil {
		return responseReq.GetResponseError(), nil
	}

	if e != nil {
		return nil, c.Config.Error
	}

	return responseReq, nil
}

func (c *ClientWA) DeleteFile(id string) types.ResponserRequest {
	// Crear request
	_, _, err := defaultRequest(http.MethodDelete, fmt.Sprintf("/%s", id), c.Config, RequestDeleteMedia)
	if err != nil {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DeleteFile request of ClientWA. error is: ", err.Error()),
		}
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DeleteFile request of ClientWA. error is: ", err.Error()),
		}
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DeleteFile request of ClientWA. error is: ", err.Error()),
		}
	}

	return types.JsonWrapperResponseRequest(b)
}

func (c *ClientWA) DeleteMessage(id string) types.ResponserRequest {
	return nil
}

// GetInfoAllNumberInWA returns information about all the phone numbers associated with the
// WhatsApp Business API client. It returns a JSON response containing an array of phone
// numbers and their associated information.
func (c *ClientWA) GetInfoAllNumberInWA() types.ResponserRequest {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", "phone_numbers"), c.Config, RequestWithQueryBusiness, QueryData{
		"access_token": CLOUD_API_ACCESS_TOKEN,
	})
	if err != nil {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllNumberInWA request of ClientWA. error is: ", err.Error()),
		}
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllNumberInWA request of ClientWA. error is: ", err.Error()),
		}
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllNumberInWA request of ClientWA. error is: ", err.Error()),
		}
	}
	fmt.Println(string(b))
	return types.JsonWrapperResponseRequest(b)
}
