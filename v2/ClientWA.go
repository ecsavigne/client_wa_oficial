//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ecsavigne/client_wa_oficial/v2/event"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/message"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
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
	MediaInfo *response.MediaInfo
}

type ClientWA struct {
	*Config `json:"config"`
}

func codeWebHook(msgByte []byte) *event.MessageWebhook {
	msg := &event.MessageWebhook{}
	json.Unmarshal(msgByte, msg)
	return msg
}

func (cl *ClientWA) initWebHookSocket() {
	// url = "wss://webhooks.savcoe-services.com/ws"
	// Conectar al servidor WebSocket
	defer func() {
		if r := recover(); r != nil {

		}
	}()

	conn, _, err := websocket.DefaultDialer.Dial(cl.Config.WebhookSocket, nil)

	var evt event.EventInterface
	if err != nil {
		switch {
		case websocket.IsUnexpectedCloseError(err):
			evt = &event.ErrorSocketConnectEvent{
				Error: response.NewError(&response.Error{
					Type:    types.TypeErrorUnexpectedClose,
					Code:    types.CodeErrorUnexpectedClose,
					Message: types.MsgErrorUnexpectedClose,
				}),
			}
		case strings.Contains(err.Error(), "tls: internal error"):
			evt = &event.ErrorSocketConnectEvent{
				Error: response.NewError(&response.Error{
					Type:    types.TypeErrorTlsInternal,
					Code:    types.CodeErrorTlsInternal,
					Message: types.MsgErrorTlsInternal,
				}),
			}
		case strings.Contains(err.Error(), "bad handshake"):
			evt = &event.ErrorSocketConnectEvent{
				Error: response.NewError(&response.Error{
					Type:    types.TypeErrorBadHandshake,
					Code:    types.CodeErrorBadHandshake,
					Message: types.MsgErrorBadHandshake,
				}),
			}
		case strings.Contains(err.Error(), "dial tcp: lookup ws"):
			evt = &event.ErrorSocketConnectEvent{
				Error: response.NewError(&response.Error{
					Type:    types.TypeErrorDialTcp,
					Code:    types.CodeErrorDialTcp,
					Message: types.MsgErrorDialTcp,
				}),
			}
		default:
			evt = &event.ErrorSocketConnectEvent{
				Error: response.NewError(&response.Error{
					Type:    types.TypeErrorUnrecognizedWebSocket,
					Code:    types.CodeErrorUnrecognizedWebSocket,
					Message: fmt.Sprintf("%s. Original error: %s", types.MsgErrorUnrecognizedWebSocket, err.Error()),
				}),
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
			evt = &event.ErrorSocketConnectEvent{
				Error: response.NewError(&response.Error{
					Type:    types.TypeErrorConectionClosedWebSocket,
					Code:    types.CodeErrorUnrecognizedWebSocket,
					Message: fmt.Sprintf("%s. Original error: %s", types.MsgErrorConectionClosedWebSocket, err.Error()),
				}),
			}
			cl.EventHandle(evt)
			break
		}

		msg := codeWebHook(message)
		switch {
		case len(msg.Entry) != 0 &&
			len(msg.Entry[0].Changes) != 0 &&
			len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
			msg.Entry[0].Changes[0].Value.Messages[0].Type != "":
			switch msg.Entry[0].Changes[0].Value.Messages[0].Type {
			case "audio":
				evt = &event.MessageAudioEvent{
					MessageWebhook: msg,
				}
			case "button":
				evt = &event.MessageButtonEvent{
					MessageWebhook: msg,
				}
			case "document":
				evt = &event.MessageDocumentEvent{
					MessageWebhook: msg,
				}
			case "text":
				evt = &event.MessageTextEvent{
					MessageWebhook: msg,
				}
			case "image":
				evt = &event.MessageImageEvent{
					MessageWebhook: msg,
				}
			case "interactive":
				evt = &event.MessageInteractiveEvent{
					MessageWebhook: msg,
				}
			case "order":
				evt = &event.MessageOrderEvent{
					MessageWebhook: msg,
				}
			case "sticker":
				evt = &event.MessageStickerEvent{
					MessageWebhook: msg,
				}
			case "system":
				evt = &event.MessageSystemEvent{
					MessageWebhook: msg,
				}
			case "video":
				evt = &event.MessageVideoEvent{
					MessageWebhook: msg,
				}
			case "reaction":
				evt = &event.MessageReactionEvent{
					MessageWebhook: msg,
				}
			case "location":
				evt = &event.MessageLocationEvent{
					MessageWebhook: msg,
				}
			case "contacts":
				evt = &event.MessageContactEvent{
					MessageWebhook: msg,
				}
			case "unknown":
				evt = &event.MessageUnknownEvent{
					MessageWebhook: msg,
				}
			default:
				cl.Config.EventHandle(message)
			}
		case len(msg.Entry) != 0 &&
			len(msg.Entry[0].Changes) != 0 &&
			len(msg.Entry[0].Changes[0].Value.Statuses) != 0:
			evt = &event.StatusMessageEvent{
				MessageWebhook: msg,
			}
		case len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
			len(msg.Entry[0].Changes[0].Value.Messages[0].Contacts) != 0:
			evt = &event.MessageContactEvent{
				MessageWebhook: msg,
			}
		default:
			cl.Config.EventHandle(message)
		}

		cl.Config.EventHandle(evt)
	}
}

// Create one Client of WhatsApp Official return *ClientWA :
// - If the EnvFilePath or Path in Config.EnvFilePath not found. ClientWA.Error = &response.Error{Type: types.TypeErrorConfig, Code: types.CodeErrorEnvNotFound, Message: types.MsgErrorEnvNotFound}
// - If occurred error in conection with WebHook Socket emit one event type: event.ErrorSocketConnectEvent
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
		c.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorBaseUrlEmpty,
			Code:    types.CodeErrorBadHandshake,
			Message: types.MsgErrorBaseUrlEmpty,
		})
		return &c
	}

	if CLOUD_API_VERSION == "" {
		c.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorApiVersionEmpty,
			Code:    types.CodeErrorApiVersionEmpty,
			Message: types.MsgErrorApiVersionEmpty,
		})
		return &c
	}

	if WA_PHONE_NUMBER_ID == "" {
		c.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorPhoneIdEmpty,
			Code:    types.CodeErrorPhoneIdEmpty,
			Message: types.MsgErrorPhoneIdEmpty,
		})
		return &c
	}

	if WA_BUSINESS_ACCOUNT_ID == "" {
		c.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorBusinessIdEmpty,
			Code:    types.CodeErrorBusinessIdEmpty,
			Message: types.MsgErrorBusinessIdEmpty,
		})
		return &c
	}

	c.path = path.Join(CLOUD_API_VERSION, WA_PHONE_NUMBER_ID)
	c.pathBusiness = path.Join(CLOUD_API_VERSION, WA_BUSINESS_ACCOUNT_ID)
	c.pathVersion = path.Join(CLOUD_API_VERSION)

	c.BaseUrl, _ = url.Parse(WA_BASE_URL)

	if c.Token == "" && CLOUD_API_ACCESS_TOKEN == "" {
		c.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorTokenEmpty,
			Code:    types.CodeErrorTokenEmpty,
			Message: types.MsgErrorTokenEmpty,
		})
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
		c.Config.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorBadRequest,
			Code:    types.CodeErrorBadRequest,
			Message: log,
		})
		return nil, c.Config.Error
	case 401:
		log := fmt.Sprintf("Error in function doRequest. Code: %d, Message: %s, MetaError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"))
		c.Config.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorUnauthorized,
			Code:    types.CodeErrorUnauthorized,
			Message: log,
		})
		return nil, c.Config.Error
	case 404:
		log := fmt.Sprintf("Error in function doRequest. Code: %d, Message: %s, MetaError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"))
		c.Config.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorUrlNotFound,
			Code:    types.CodeErrorUrlNotFound,
			Message: log,
		})
		return nil, c.Config.Error
	}

	return res, nil
}

func (c *ClientWA) doRequest(req *http.Request) (response.ResponserRequest, error) {
	res, e := doRequest(req, c)

	var responser response.ResponserRequest
	if e != nil {
		responser = response.NewError(&response.Error{
			Type:    types.TypeErrorInRequest,
			Code:    types.CodeErrorInRequest,
			Message: fmt.Sprintf("Type: %s. Error is: %s", types.MsgErrorInRequest, e.Error()),
		})
		return responser, c.Config.Error
	}

	bodyResponse, e := io.ReadAll(res.Body)
	if e != nil {
		log := fmt.Sprintf("Error in function doRequest of ClientWA when reading response body. Error is: %s", e.Error())
		responser = response.NewError(&response.Error{
			Type:    types.TypeErrorUnrecognized,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Type: %s. Error is: %s", types.MsgErrorUnrecognized, e.Error()),
		})
		c.Config.Error = fmt.Errorf("%s", log)
		return nil, c.Config.Error
	}

	defer res.Body.Close()

	return response.JsonWrapperResponseRequest(bodyResponse), nil
}

func (c *ClientWA) makeRequest(methoth string, ePoint string, msg message.Messager) (response.ResponserRequest, error) {
	if msg.GetMessageLink() != "" || msg.GetFileHeader() != nil {
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
		responseReq response.ResponserRequest
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

func validTypeMsg(msg message.Messager, msgType string) bool {
	return msg.GetType() == msgType
}

// SendTemplate sends a template message. It validates the message type to ensure it is a template.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendTemplate(m message.Messager) response.ResponserRequest {

	if !validTypeMsg(m, types.MessageTypeTemplate) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeTemplate, m.GetType()),
		})
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendTemplate request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendTextMessage sends a text message. It validates the message type to ensure it is a text.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendTextMessage(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeText) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeText, m.GetType()),
		})
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendText request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendReaction sends a reaction message. It validates the message type to ensure it is a reaction.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendReaction(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeReaction) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeReaction, m.GetType()),
		})
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendReaction request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendInteractiveList sends an interactive list message. It validates the message type to ensure it is an interactive of type list.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendInteractiveList(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(types.InteractiveTypeList) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeList),
			})
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveList request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendInteractiveButtonResponse sends an interactive button response message. It validates the message type to ensure it is an
// interactive of type button response. If the message type is incorrect, it returns an error. Otherwise, it makes a request to
// send the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendInteractiveButtonResponse(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m.(*message.MessageInteractive), types.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(types.InteractiveTypeButtonResponse) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeButtonResponse),
			})
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveButtonResponse request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendInteractiveButtonUrl sends an interactive button URL message. It validates the message type to ensure it is an
// interactive of type button URL. If the message type is incorrect, it returns an error. Otherwise, it makes a request
// to send the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendInteractiveButtonUrl(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m.(*message.MessageInteractive), types.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(types.InteractiveTypeButtonUrl) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeButtonUrl),
			})
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveButtonUrl request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendInteractiveMsgProcess sends an interactive process message. It validates the message type to ensure it is an
// interactive of type process. If the message type is incorrect, it returns an error. Otherwise, it makes a request
// to send the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendInteractiveMsgProcess(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m.(*message.MessageInteractive), types.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(types.InteractiveTypeProcess) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeProcess),
			})
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveMsgProcess request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendInteractiveOneProduct sends an interactive message of type product. It validates the message type to ensure it is an
// interactive of type product. If the message type is incorrect, it returns an error. Otherwise, it makes a request
// to send the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendInteractiveOneProduct(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m.(*message.MessageInteractive), types.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(types.InteractiveTypeProduct) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeProduct),
			})
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveOneProduct request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendInteractiveMultiProduct sends an interactive message of type multi product. It validates the message type to
// ensure it is an interactive of type multi product. If the message type is incorrect, it returns an error. Otherwise,
// it makes a request to send the message. If the request is successful, the response is returned; otherwise, an error
// is returned.
func (c *ClientWA) sendInteractiveMultiProduct(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m.(*message.MessageInteractive), types.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(types.InteractiveTypeMultiProduct) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeMultiProduct),
			})
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveMultiProduct request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendInteractiveCatalog sends an interactive catalog message. It validates the message type
// to ensure it is an interactive of type catalog. If the message type is incorrect, it returns
// an error. Otherwise, it makes a request to send the message. If the request is successful,
// the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendInteractiveCatalog(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m.(*message.MessageInteractive), types.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(types.InteractiveTypeCatalog) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", types.InteractiveTypeCatalog),
			})
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendInteractiveCatalog request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendResponseMsg sends a response message. It validates the message type to ensure it is a response
// message. If the message type is incorrect, it returns an error. Otherwise, it makes a request to send
// the message. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendResponseMsg(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, "response") {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", "response", m.GetType()),
		})
	} else {
		if !m.IsTypeResponse() {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message type response not expect, type: '%s'", m.(*message.MessageResponse).MessagerKernel.Type),
			})
		}
	}

	if m.GetMessageLink() != "" || m.GetFileHeader() != nil {
		r := c.UploadFile(m, message.AUDIO)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendResponseMsg request of ClientWA. error is: ", e.Error()),
		})
	}

	if resp.GetType() == response.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

func (c *ClientWA) validLinAndId(m message.Messager) response.ResponserRequest {
	if m.GetMessageLink() != "" && m.GetMessageId() != "" {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: "Expect Message.id or Message.link, but not both",
		})
	}

	return nil
}

// SendAudioMessage sends an audio message. It validates the message type to ensure it is an audio.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendAudioMessage(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeAudio) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeAudio, m.GetType()),
		})
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" || m.GetFileHeader() != nil {
		r := c.UploadFile(m, message.AUDIO)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendAudio request of ClientWA. error is: ", e.Error()),
		})
	}

	if resp.GetType() == response.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendImageMessage sends an image message. It validates the message type to ensure it is an image.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendImageMessage(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeImage) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeImage, m.GetType()),
		})
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" || m.GetFileHeader() != nil {
		r := c.UploadFile(m, message.IMAGE)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendImage request of ClientWA. error is: ", e.Error()),
		})
	}

	if resp.GetType() == response.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendVideoMessage sends a video message. It validates the message type to ensure it is a video.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendVideoMessage(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeVideo) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeVideo, m.GetType()),
		})
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" || m.GetFileHeader() != nil {
		r := c.UploadFile(m, message.VIDEO)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendVideo request of ClientWA. error is: ", e.Error()),
		})
	}

	if resp.GetType() == response.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendDocumentMessage sends a document message. It validates the message type to ensure it is a document.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendDocumentMessage(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeDocument) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeDocument, m.GetType()),
		})
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" || m.GetFileHeader() != nil {
		r := c.UploadFile(m, message.DOCUMENT)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendDocument request of ClientWA. error is: ", e.Error()),
		})
	}

	if resp.GetType() == response.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendStickerMessage sends a sticker message. It validates the message type to ensure it is a sticker. If the message type is
// incorrect, it returns an error. Otherwise, it makes a request to send the message. If the request is successful, the response
// is returned; otherwise, an error is returned.
func (c *ClientWA) sendStickerMessage(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeSticker) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeSticker, m.GetType()),
		})
	} else if r := c.validLinAndId(m); r != nil {
		return r
	}

	if m.GetMessageLink() != "" || m.GetFileHeader() != nil {
		r := c.UploadFile(m, message.STICKER)
		if r != nil && r.GetResponseError() != nil {
			return r.GetResponseError()
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendSticker request of ClientWA. error is: ", e.Error()),
		})
	}

	if resp.GetType() == response.ResponseSuccess {
		resp.GetResponseSuccess().MediaInfo = c.MediaInfo
		c.resetMessageInfo()
	}

	return resp
}

// SendLocationMessage sends a location message. It validates the message type
// to ensure it is a location message. If the message type is incorrect, it
// returns an error. Otherwise, it makes a request to send the message. If the
// request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendLocationMessage(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeLocation) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeLocation, m.GetType()),
		})
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendLocationMessage request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// SendContactMessage sends a contact message. It first validates the message type to ensure it is a contact message.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendContactMessage(m message.Messager) response.ResponserRequest {
	if !validTypeMsg(m, types.MessageTypeContact) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeContact, m.GetType()),
		})
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		if ok := e.(*response.Error); ok != nil {
			return e.(*response.Error)
		}

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error en SendContactMessage request of ClientWA. error is: ", e.Error()),
		})
	}

	return resp
}

// func (c *ClientWA) sendTemplate(m message.Messager) response.ResponserRequest {
// 	return nil
// 	// switch m.() {
// 	// case types.MessageTypeTemplate:
// 	// }
// }

func (c *ClientWA) sendInteractive(m message.Messager) response.ResponserRequest {
	interactive := m.GetInteractiveMessage()
	if interactive != nil {
		switch interactive.GetInteractiveProto().Type {
		case types.InteractiveTypeList:
			return c.sendInteractiveList(m)
		case types.InteractiveTypeButtonResponse:
			return c.sendInteractiveButtonResponse(m)
		case types.InteractiveTypeProduct:
			c.sendInteractiveOneProduct(m)
		case types.InteractiveTypeMultiProduct:
			return c.sendInteractiveMultiProduct(m)
		case types.InteractiveTypeProcess:
			return c.sendInteractiveMsgProcess(m)
		case types.InteractiveTypeCatalog:
			return c.sendInteractiveCatalog(m)
		case types.InteractiveTypeButtonUrl:
			return c.sendInteractiveButtonUrl(m)
		}
	}
	return response.NewError(&response.Error{
		Type:    response.ResponseError,
		Code:    401,
		Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeInteractive, m.GetType()),
	})
}

func (c *ClientWA) SendMessage(m message.Messager) response.ResponserRequest {
	switch m.GetType() {
	case types.MessageTypeAudio:
		return c.sendAudioMessage(m)
	case types.MessageTypeContact:
		return c.sendContactMessage(m)
	case types.MessageTypeDocument:
		return c.sendDocumentMessage(m)
	case types.MessageTypeImage:
		return c.sendImageMessage(m)
	case types.MessageTypeInteractive:
		return c.sendInteractive(m)
	case types.MessageTypeLocation:
		return c.sendLocationMessage(m)
	case types.MessageTypeReaction:
		return c.sendReaction(m)
	case types.MessageTypeResponse:
		return c.sendResponseMsg(m)
	case types.MessageTypeSticker:
		return c.sendStickerMessage(m)
	case types.MessageTypeTemplate:
		return c.sendTemplate(m)
	case types.MessageTypeText:
		return c.sendTextMessage(m)
	case types.MessageTypeVideo:
		return c.sendVideoMessage(m)
	default:
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message not recognized. Message.type expect '%v'", []string{"text", "audio", "image", "video", "document", "sticker", "location", "contact", "template", "interactive", "reaction"}),
		})
	}
}

// UploadFile uploads a file to the server. It validates the message type to ensure it is either a text, audio, image, video,
// document, or sticker message. If the message type is incorrect, it returns an error. Otherwise, it makes a request to upload
// the file. If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) UploadFile(m message.Messager, mt message.MediaType) response.ResponserRequest {
	resp, e := c.makeRequest(http.MethodPost, "/media", m)
	if e != nil {
		msgError := fmt.Sprintln("Error in UploadFile request of ClientWA. error is: ", e.Error())
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: msgError,
		})

	} else if resp.GetType() == response.ResponseError {
		return resp.GetResponseError()
	}

	id := ""

	if resp.GetType() == response.ResponseMediaInfo {
		c.Config.MediaInfo = resp.GetResponseMediaInfo()
		id = c.Config.MediaInfo.ID
	}

	switch mt {
	case message.AUDIO, message.IMAGE, message.VIDEO, message.DOCUMENT, message.STICKER:
		m.SetId(id)
		m.SetLink("")
	}

	return resp
}

// DownloadFile downloads a file using its unique identifier. It first retrieves the file information
// and checks if the response indicates an error or contains media information. If media information
// is available, it sends a request to download the file and saves it to the specified path with the
// given name whithout extension. If an error occurs during any step, it returns an error. Otherwise, it completes the
// if param save == true, the file is saved to the specified path with the given name else it returns the response of the request.
// download process successfully.
func (c *ClientWA) DownloadFile(id, path, nameFile string, save ...bool) (*http.Response, context.CancelFunc, error) {
	saveValue := true
	if len(save) == 1 {
		saveValue = false
	}
	responseReq, e := c.getFileInfo(id)
	if e != nil {
		return nil, nil, c.Config.Error
	}

	if responseReq.IsType(response.ResponseError) {
		return nil, nil, responseReq.GetResponseError()
	}

	if responseReq.IsType(response.ResponseMediaInfo) {
		// Get binaryFile
		mInfo := responseReq.GetResponseMediaInfo()
		_, cancel, _ := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", mInfo.Url), c.Config, RequestChangeUrlFull)

		if c.Config.Error != nil {
			return nil, nil, c.Config.Error
		} else {
			// Save binaryFile in path
			res, e := doRequest(c.request, c)
			ext := strings.Split(res.Header.Get("Content-Disposition"), ".")
			if e != nil {
				return nil, nil, e
			} else {
				if !saveValue {
					var resTemp http.Response = *res
					return &resTemp, cancel, e
				}
				file, e := os.Create(fmt.Sprintf("%s%s.%s", path, nameFile, ext[len(ext)-1]))
				if e != nil {
					return nil, nil, e
				}
				defer file.Close()

				_, e = io.Copy(file, res.Body)
				if e != nil {
					return nil, nil, response.NewError(&response.Error{
						Type:    response.ResponseError,
						Code:    types.CodeErrorUnrecognized,
						Message: fmt.Sprintln("Error in DownloadFile request of ClientWA. error is: ", e.Error()),
					})
				}
				defer res.Body.Close()
				defer cancel()
			}
		}
	}
	return nil, nil, nil
}

func (c *ClientWA) getFileInfo(id string) (response.ResponserRequest, error) {
	// Crear request
	defaultRequest(http.MethodGet, fmt.Sprintf("/%s", id), c.Config, RequestGetMessageInfo)

	var (
		responseReq response.ResponserRequest
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

func (c *ClientWA) DeleteFile(id string) response.ResponserRequest {
	// Crear request
	_, _, err := defaultRequest(http.MethodDelete, fmt.Sprintf("/%s", id), c.Config, RequestDeleteMedia)
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DeleteFile request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DeleteFile request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DeleteFile request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b)
}

func (c *ClientWA) DeleteMessage(id string) response.ResponserRequest {
	return nil
}

// GetInfoAllNumberInWA returns information about all the phone numbers associated with the
// WhatsApp Business API client. It returns a JSON response containing an array of phone
// numbers and their associated information.
func (c *ClientWA) GetInfoAllNumberInWA() response.ResponserRequest {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", "phone_numbers"), c.Config, RequestWithQueryBusiness, QueryData{
		"access_token": CLOUD_API_ACCESS_TOKEN,
	})
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllNumberInWA request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllNumberInWA request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllNumberInWA request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b)
}
