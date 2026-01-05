//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path"
	"slices"
	str "strings"
	"sync"
	"time"

	"github.com/ecsavigne/client_wa_oficial/v2/event"
	evt_types "github.com/ecsavigne/client_wa_oficial/v2/event/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/message"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
	"golang.org/x/net/http2"

	"github.com/gorilla/websocket"
)

func (c CLIENT_TYPE) String() string {
	return string(c)
}

const (
	CLIENT_WHATSAPP CLIENT_TYPE = "whatsapp"
	CLIENT_FACEBOOK CLIENT_TYPE = "facebook"
)

type (
	CLIENT_TYPE string
	clientHttp  struct {
		*http.Client
		BaseUrl *url.URL `json:"base_url"`
	}

	ClientWA struct {
		*Config    `json:"config"`
		typeClient CLIENT_TYPE
	}

	InfoContact struct {
		ContactPhone string
		RecipientID  string
		IsOnWhats    bool
		IsError      bool
		Error        error
		MsgError     string
	}

	chanDataWaba struct {
		WabaInfo    response.WabaInfo
		PhoneInfo   *response.PhoneInfo
		ExistNumber bool
	}

	pair struct {
		Phone   string
		Channel chan InfoContact
	}
)

var (
	infoContacts map[string]pair
)

func codeWebHook(msgByte []byte) *event.Components {
	msg := &event.Components{}
	json.Unmarshal(msgByte, msg)
	return msg
}

// messageIsForMe && Field
// field = template_category_update, message_template_status_update, messages
func (cl *ClientWA) messageIsForMe(component *event.Components) (isForme bool, typeNotification evt_types.TYPE_NOTIFICATION_WEBHOOK) {
	field := component.Entry[0].Changes[0].Field
	phoneNumberID := ""
	if component.Entry[0].Changes[0].Value.Metadata != nil {
		phoneNumberID = component.Entry[0].Changes[0].Value.Metadata.PhoneNumberID
	}
	wabaID := component.Entry[0].ID

	if phoneNumberID == cl.Config.getphoneID() && field == "messages" {
		isForme = true

		return isForme, evt_types.ParseTypeNotificationWebhook(field)
	}

	if wabaID == cl.Config.getWabaID() && (field == "template_category_update" || field == "message_template_status_update") {
		isForme = true
	}

	return isForme, evt_types.ParseTypeNotificationWebhook(field)
}

func (cl ClientWA) GetType() string {
	return cl.typeClient.String()
}

func (cl *ClientWA) initWebHookSocket() {
	// url = "wss://webhooks.savcoe-services.com/ws"
	// Conectar al servidor WebSocket
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	var evt event.EventInterface

	conn, _, err := websocket.DefaultDialer.Dial(cl.Config.WebhookSocket, nil)
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
		case str.Contains(err.Error(), "tls: internal error"):
			evt = &event.ErrorSocketConnectEvent{
				Error: response.NewError(&response.Error{
					Type:    types.TypeErrorTlsInternal,
					Code:    types.CodeErrorTlsInternal,
					Message: types.MsgErrorTlsInternal,
				}),
			}
		case str.Contains(err.Error(), "bad handshake"):
			evt = &event.ErrorSocketConnectEvent{
				Error: response.NewError(&response.Error{
					Type:    types.TypeErrorBadHandshake,
					Code:    types.CodeErrorBadHandshake,
					Message: types.MsgErrorBadHandshake,
				}),
			}
		case str.Contains(err.Error(), "dial tcp: lookup ws"):
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
		isForme, _ := cl.messageIsForMe(msg)
		if !isForme {
			continue
		}

		switch {
		case len(msg.Entry) != 0 &&
			len(msg.Entry[0].Changes) != 0 &&
			len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
			msg.Entry[0].Changes[0].Value.Messages[0].Type != "":
			switch msg.Entry[0].Changes[0].Value.Messages[0].Type {
			case "audio":
				evt = &event.MessageAudioEvent{
					Components: msg,
				}
			case "button":
				evt = &event.MessageButtonEvent{
					Components: msg,
				}
			case "document":
				evt = &event.MessageDocumentEvent{
					Components: msg,
				}
			case "text":
				evt = &event.MessageTextEvent{
					Components: msg,
				}
			case "image":
				evt = &event.MessageImageEvent{
					Components: msg,
				}
			case "interactive":
				evt = &event.MessageInteractiveEvent{
					Components: msg,
				}
			case "order":
				evt = &event.MessageOrderEvent{
					Components: msg,
				}
			case "sticker":
				evt = &event.MessageStickerEvent{
					Components: msg,
				}
			case "system":
				evt = &event.MessageSystemEvent{
					Components: msg,
				}
			case "video":
				evt = &event.MessageVideoEvent{
					Components: msg,
				}
			case "reaction":
				evt = &event.MessageReactionEvent{
					Components: msg,
				}
			case "location":
				evt = &event.MessageLocationEvent{
					Components: msg,
				}
			case "contacts":
				evt = &event.MessageContactEvent{
					Components: msg,
				}
			case "unknown":
				evt = &event.MessageUnknownEvent{
					Components: msg,
				}
			default:
				cl.Config.EventHandle(message)
			}
		case len(msg.Entry) != 0 &&
			len(msg.Entry[0].Changes) != 0 &&
			len(msg.Entry[0].Changes[0].Value.Statuses) != 0:
			evt = &event.StatusMessageEvent{
				Components: msg,
			}

			recipientId := msg.Entry[0].Changes[0].Value.Statuses[0].RecipientID
			id := msg.Entry[0].Changes[0].Value.Statuses[0].ID

			mu := sync.Mutex{}
			mu.Lock()
			pair, ok := infoContacts[id]
			mu.Unlock()

			if msg.Entry[0].Changes[0].Value.Statuses[0].Status == "failed" &&
				msg.Entry[0].Changes[0].Value.Statuses[0].Errors[0].Message == "Message undeliverable" {
				if ok {
					pair.Channel <- InfoContact{
						ContactPhone: pair.Phone,
						RecipientID:  recipientId,
						IsOnWhats:    false,
					}
				}

			} else {
				if ok {
					pair.Channel <- InfoContact{
						ContactPhone: pair.Phone,
						RecipientID:  recipientId,
						IsOnWhats:    true,
					}
				}
			}
		case len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
			len(msg.Entry[0].Changes[0].Value.Messages[0].Contacts) != 0:
			evt = &event.MessageContactEvent{
				Components: msg,
			}
		default:
			cl.Config.EventHandle(message)
		}

		cl.Config.EventHandle(evt)
	}
}

func defaultConfig() *Config {
	return &Config{
		WebhookSocket:     "",
		EventHandle:       nil,
		wabaID:            "",
		businessID:        "",
		phoneID:           "",
		cLOUD_API_VERSION: "v24.0",
	}
}

func NewClientWA(opts ...Options) *ClientWA {
	c := *defaultConfig()

	cl := &ClientWA{
		typeClient: CLIENT_WHATSAPP,
		Config:     &c,
	}

	for _, opt := range opts {
		opt(&c)
	}

	err := setEnv(&c)
	if err != nil {
		cl.Error = err
		return cl
	}

	cl.Config = newConfig(c)
	if cl.Error != nil {
		return cl
	}

	if c.WebhookSocket != "" && c.EventHandle != nil {
		go cl.initWebHookSocket()
	}

	return cl
}

func createClientHttp2(timeOut int) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12},
	}
	http2.ConfigureTransport(tr)

	return &http.Client{
		Timeout:   time.Duration(timeOut) * time.Second,
		Transport: tr,
	}
}

func newConfig(c Config) *Config {
	c.Error = nil
	if c.wA_BASE_URL == "" {
		c.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorBaseUrlEmpty,
			Code:    types.CodeErrorBadHandshake,
			Message: types.MsgErrorBaseUrlEmpty,
		})
		return &c
	} else {
		c.BaseUrl, _ = url.Parse(c.wA_BASE_URL)
	}

	if c.cLOUD_API_VERSION == "" {
		c.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorApiVersionEmpty,
			Code:    types.CodeErrorApiVersionEmpty,
			Message: types.MsgErrorApiVersionEmpty,
		})
		return &c
	} else {
		c.pathVersion = path.Join(c.cLOUD_API_VERSION)
	}

	// if c.getphoneID() == "" {
	// 	c.pathPhone = path.Join(c.pathVersion, c.phoneID)
	// c.Error = response.NewError(&response.Error{
	// 	Type:    types.TypeErrorPhoneIdEmpty,
	// 	Code:    types.CodeErrorPhoneIdEmpty,
	// 	Message: types.MsgErrorPhoneIdEmpty,
	// })
	// return &c
	// } else {
	c.pathPhone = path.Join(c.getVersion(), c.phoneID)
	// }

	// if c.getWabaID() == "" {
	// 	c.pathWaba = path.Join(c.pathVersion, c.wabaID)
	// c.Error = response.NewError(&response.Error{
	// 	Type:    types.TypeErrorBusinessIdEmpty,
	// 	Code:    types.CodeErrorBusinessIdEmpty,
	// 	Message: types.MsgErrorBusinessIdEmpty,
	// })
	// return &c
	// } else {
	c.pathWaba = path.Join(c.getVersion(), c.wabaID)
	// }

	// if c.getBusinessID() == "" {
	// 	c.pathBusiness = path.Join(c.pathVersion, c.businessID)
	// c.Error = response.NewError(&response.Error{
	// 	Type:    types.TypeErrorBusinessIdEmpty,
	// 	Code:    types.CodeErrorBusinessIdEmpty,
	// 	Message: types.MsgErrorBusinessIdEmpty,
	// })
	// return &c
	// } else {
	c.pathBusiness = path.Join(c.getVersion(), c.businessID)
	// }

	if c.Token == "" && c.cLOUD_API_ACCESS_TOKEN == "" {
		c.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorTokenEmpty,
			Code:    types.CodeErrorTokenEmpty,
			Message: types.MsgErrorTokenEmpty,
		})
		return &c
	}

	if c.Client == nil {
		c.Client = createClientHttp2(30)
	}

	return &c
}

func (c *ClientWA) resetMessageInfo() {
	if c.MediaInfo != nil {
		c.MediaInfo = nil
	}
}

func doRequest(req *http.Request, c *ClientWA) (*http.Response, error) {
	c.clientHttp.Client = createClientHttp2(30)
	res, err := c.clientHttp.Do(req)
	if err != nil {
		log := fmt.Sprintf("Error in function doRequest when send HTTP request to server with Do. Error is: %s", err.Error())
		c.Config.Error = fmt.Errorf("%s", log)
		return nil, c.Config.Error
	}

	switch res.StatusCode {
	case 400:
		if e, ok := response.GetResponseRequest(res.Body, "doRequest", "400").(*response.Error); ok {
			c.Config.Error = e
			return nil, e
		}
		log := fmt.Sprintf("Error in function doRequest. Code: %d, Message: %s, MetaError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"))
		c.Config.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorBadRequest,
			Code:    types.CodeErrorBadRequest,
			Message: log,
		})
		return nil, c.Config.Error
	case 401:
		if e, ok := response.GetResponseRequest(res.Body, "doRequest", "401").(*response.Error); ok {
			c.Config.Error = e
			return nil, e
		}
		log := fmt.Sprintf("Error in function doRequest. Code: %d, Message: %s, MetaError: %s.", res.StatusCode, res.Status, res.Header.Get("Www-Authenticate"))
		c.Config.Error = response.NewError(&response.Error{
			Type:    types.TypeErrorUnauthorized,
			Code:    types.CodeErrorUnauthorized,
			Message: log,
		})
		return nil, c.Config.Error
	case 404:
		if e, ok := response.GetResponseRequest(res.Body, "doRequest", "404").(*response.Error); ok {
			c.Config.Error = e
			return nil, e
		}
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

func (c *ClientWA) doRequest(req *http.Request) (response.Responser, error) {
	res, e := doRequest(req, c)

	var responser response.Responser
	if e != nil {
		responser = response.NewError(&response.Error{
			Type:    types.TypeErrorInRequest,
			Code:    types.CodeErrorInRequest,
			Message: fmt.Sprintf("Type: %s. Error is: %s", types.MsgErrorInRequest, e.Error()),
		})
		return responser, c.Config.Error
	}
	defer res.Body.Close()

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

	return response.JsonWrapperResponseRequest(bodyResponse), nil
}

func (c *ClientWA) makeRequest(methoth string, ePoint string, msg message.Messager) (response.Responser, error) {
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
		responseReq response.Responser
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
func (c *ClientWA) sendTemplate(m message.Messager) response.Responser {

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
func (c *ClientWA) sendTextMessage(m message.Messager) response.Responser {
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
func (c *ClientWA) sendReaction(m message.Messager) response.Responser {
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
func (c *ClientWA) sendInteractiveList(m message.Messager) response.Responser {
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
func (c *ClientWA) sendInteractiveButtonResponse(m message.Messager) response.Responser {
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
func (c *ClientWA) sendInteractiveButtonUrl(m message.Messager) response.Responser {
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
func (c *ClientWA) sendInteractiveMsgProcess(m message.Messager) response.Responser {
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
func (c *ClientWA) sendInteractiveOneProduct(m message.Messager) response.Responser {
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
func (c *ClientWA) sendInteractiveMultiProduct(m message.Messager) response.Responser {
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
func (c *ClientWA) sendInteractiveCatalog(m message.Messager) response.Responser {
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
func (c *ClientWA) sendResponseMsg(m message.Messager) response.Responser {
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

func (c *ClientWA) validLinAndId(m message.Messager) response.Responser {
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
func (c *ClientWA) sendAudioMessage(m message.Messager) response.Responser {
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
func (c *ClientWA) sendImageMessage(m message.Messager) response.Responser {
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
func (c *ClientWA) sendVideoMessage(m message.Messager) response.Responser {
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
func (c *ClientWA) sendDocumentMessage(m message.Messager) response.Responser {
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
func (c *ClientWA) sendStickerMessage(m message.Messager) response.Responser {
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
func (c *ClientWA) sendLocationMessage(m message.Messager) response.Responser {
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
func (c *ClientWA) sendContactMessage(m message.Messager) response.Responser {
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

// func (c *ClientWA) sendTemplate(m message.Messager) response.Responser {
// 	return nil
// 	// switch m.() {
// 	// case types.MessageTypeTemplate:
// 	// }
// }

func (c *ClientWA) sendInteractive(m message.Messager) response.Responser {
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

func (c *ClientWA) SendMessage(m message.Messager) response.Responser {
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
func (c *ClientWA) UploadFile(m message.Messager, mt message.MediaType) response.Responser {
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
			ext := str.Split(res.Header.Get("Content-Disposition"), ".")
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

func (c *ClientWA) getFileInfo(id string) (response.Responser, error) {
	// Crear request
	defaultRequest(http.MethodGet, fmt.Sprintf("/%s", id), c.Config, RequestGetMessageInfo)

	var (
		responseReq response.Responser
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

func (c *ClientWA) DeleteFile(id string) response.Responser {
	// Crear request
	_, _, err := defaultRequest(http.MethodDelete, fmt.Sprintf("/%s", id), c.Config, RequestDeleteMedia)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DeleteFile request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
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

func (c *ClientWA) DeleteMessage(id string) response.Responser {
	return nil
}

// GetInfoAllNumberInWaba retrieves information about all phone numbers associated with a given
// WhatsApp Business Account (specified by waba_id).
func (c *ClientWA) GetInfoAllNumberInWaba(waba_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s/%s", waba_id, "phone_numbers"), c.Config, RequestWithVersion, nil)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllNumberInWA request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
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

func (c *ClientWA) GetNumberInfo(phone_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", phone_id), c.Config, RequestWithVersion, nil)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b)
}

// GetInfoAllNumberInWA returns information about all the Whatsapp Business Account phone associated with the
// WhatsApp Business API client. It returns a JSON response containing an array of phone
// numbers and their associated information.
func (c *ClientWA) GetOwnedWaba(portafolio_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s/%s", portafolio_id, "owned_whatsapp_business_accounts"), c.Config, RequestWithVersion)
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllWaba request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllWaba request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoAllWaba request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b)
}

func validKeyInMapInFunc(data map[string]any, key []string, funcName string) response.Responser {
	if data == nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Error in %s. error is: params data is required and is cannot absent or nil", funcName),
		})
	}

	keysData := slices.Collect(maps.Keys(data))
	for _, v := range slices.All(key) {
		if !slices.Contains(keysData, v) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    types.CodeErrorUnrecognized,
				Message: fmt.Sprintf("Error in %s. error is: %s required and cannot be absent", funcName, v),
			})
		} else {
			if data[v] == "" {
				return response.NewError(&response.Error{
					Type:    response.ResponseError,
					Code:    types.CodeErrorUnrecognized,
					Message: fmt.Sprintf("Error in %s. error is: %s required and cannot be empty", funcName, v),
				})
			}
		}
	}

	return nil
}

func (c *ClientWA) RegisterNumberInWaba(data map[string]any) response.Responser {
	if err := validKeyInMapInFunc(data, []string{"cc", "phone_number", "verified_name"}, "RegisterNumberInWaba"); err != nil {
		return err
	}

	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%s/%s", c.GetWabaId(), "phone_numbers"), c.Config, RequestWithVersion, data)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in RegisterNumberInWaba request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in RegisterNumberInWaba request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "RegisterNumberInWaba", "ClientWA")
}

func (c *ClientWA) GetVerificationCode(data map[string]any) response.Responser {
	if err := validKeyInMapInFunc(data, []string{"code_method", "language"}, "GetverificationCode"); err != nil {
		return err
	}

	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%s", "request_code"), c.Config, RequestWithQueryPhone, QueryData(data))
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetverificationCode request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetverificationCode request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "GetverificationCode", "ClientWA")
}

func (c *ClientWA) VerifyCode(data map[string]any) response.Responser {
	if err := validKeyInMapInFunc(data, []string{"code"}, "VerifyCode"); err != nil {
		return err
	}

	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%s", "verify_code"), c.Config, RequestWithQueryPhone, QueryData(data))
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in VerifyCode request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in VerifyCode request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "VerifyCode", "ClientWA")
}

func (c *ClientWA) RegisterForUseApi() response.Responser {
	data := map[string]any{
		"messaging_product": "whatsapp",
		"pin":               "123456",
	}

	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%s", "register"), c.Config, data)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in RegisterForUseApi request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in RegisterForUseApi request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "RegisterForUseApi", "ClientWA")
}

func (c *ClientWA) setWabaID(waba_id string) {
	c.Config.setWabaID(waba_id)
}

func (c *ClientWA) SetPhoneId(phone_id string) {
	c.Config.setPhoneID(phone_id)
}

func (c ClientWA) GetWabaId() string {
	return c.getWabaID()
}

func (c ClientWA) GetPhoneId() string {
	return c.getphoneID()
}

// GetWabaInfo makes a GET call to the WhatsApp API to get info about
// a WhatsApp Business Account by its ID.
//
// It returns a ResponserRequest with the response data in the Body
// field, which is a JSON object with the following structure:
//
//	{
//	    "id": "1234567890",
//	    "name": "Example Inc",
//	    "quality_rating": "GREEN",
//	    "status": "ACTIVE",
//	    "created_time": "1562344562",
//	    "updated_time": "1562344562"
//	}
//
// If the request fails, it returns a ResponserRequest with the error
func (c *ClientWA) GetWabaInfo(waba_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", waba_id), c.Config, RequestWithVersion, nil)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b)
}

// GetBusinessInfo retrieves detailed information about a business using its business ID.
// It performs a GET request to fetch fields such as id, name, extended updated time,
// link, two-factor type, is_hidden status, payment account ID, verification status,
// updated time, and created time. If an error occurs during the request, a response
// error is returned. The response is wrapped and returned as a ResponserRequest.

func (c *ClientWA) GetBusinessInfo(business_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", business_id), c.Config, RequestWithQueryVersion, QueryData{"fields": "id,name,extended_updated_time,link,two_factor_type,is_hidden,payment_account_id,verification_status,updated_time,created_time"})
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetNumberInfo request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b)
}

func getPhoneNumber(phone string) string {
	return str.NewReplacer("+", "", "-", "", " ", "").Replace(phone)
}

// GetInfoPhoneOfWaba returns info of a phone in a waba.
//
// Example:
// resp := c.GetInfoPhoneOfWaba("+573123456789", "1234567890123")
func (c *ClientWA) GetInfoPhoneOfWaba(phoneNumber, waba_id string) response.Responser {
	resp := c.GetInfoAllNumberInWaba(waba_id)

	if nums := resp.GetResponsePhonesWA(); nums != nil {
		for _, phoneInfo := range nums.Data {
			if getPhoneNumber(phoneInfo.DisplayPhoneNumber) == phoneNumber {
				return response.NewPhone(&response.Phone{
					PhoneInfo: &phoneInfo,
				})
			}
		}
	}

	return response.NewError(&response.Error{
		Type:    response.ResponseError,
		Code:    types.CodeErrorUnrecognized,
		Message: fmt.Sprintf("Phone number: %s not found in waba-id: %s in Meta", phoneNumber, waba_id),
	})
}

func workerFindWabaId(ctx context.Context, wg *sync.WaitGroup, cl *ClientWA, dIn <-chan chanDataWaba, dOut chan<- chanDataWaba, phone_number string) {
	defer wg.Done()

	for data := range dIn {
		resp := cl.GetInfoAllNumberInWaba(data.WabaInfo.ID)
		if nums := resp.GetResponsePhonesWA(); nums != nil {
			for _, phoneInfo := range nums.Data {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if getPhoneNumber(phoneInfo.DisplayPhoneNumber) == phone_number {
					data.ExistNumber = true
					data.PhoneInfo = &phoneInfo
					dOut <- data
					return
				}
			}
		}
	}
}

// FindWabaId: Find the WabaId and PhoneInfo of a phone_number associated to a portafolio_id in Meta.
// If not found, return nil, nil.
func (c *ClientWA) FindWabaId(portafolio_id, phone_number string) (*response.WabaInfo, *response.PhoneInfo) {
	wabas := c.GetOwnedWaba(portafolio_id)
	arrWabaInfo := wabas.GetResponseWaba().Data
	cant := len(arrWabaInfo)
	arrFuncCancel := make([]context.CancelFunc, 0)
	var wg sync.WaitGroup

	// create channel
	dIn := make(chan chanDataWaba, cant)
	dOut := make(chan chanDataWaba, cant)

	// create workers
	cantWorker := cant / 5
	if cant%5 != 0 {
		cantWorker++
	}

	for range cantWorker {
		ctx, fCancel := context.WithCancel(context.Background())
		arrFuncCancel = append(arrFuncCancel, fCancel)

		wg.Add(1)
		go workerFindWabaId(ctx, &wg, c, dIn, dOut, phone_number)
	}

	// send data to workers
	go func() {
		// close channel
		defer close(dIn)

		for _, wabaInfo := range arrWabaInfo {
			dIn <- chanDataWaba{
				WabaInfo: wabaInfo,
			}
		}
	}()

	// close channel out
	go func() {
		wg.Wait()
		close(dOut)
	}()

	for data := range dOut {
		if data.ExistNumber {
			// set WaBusinessAccountId
			c.setWabaID(data.WabaInfo.ID)

			// Cancel context
			for _, fCancel := range arrFuncCancel {
				fCancel()
			}
			c.Config.Error = nil

			return &data.WabaInfo, data.PhoneInfo
		}
	}

	return nil, nil
}

func registerContact(idMsg, phone string) chan InfoContact {
	ch := make(chan InfoContact, 1)

	p := pair{
		Phone:   phone,
		Channel: ch,
	}

	mu := &sync.Mutex{}
	mu.Lock()
	infoContacts[idMsg] = p
	mu.Unlock()

	return ch
}

func getMsgPing(cl *ClientWA, num string) (msgID string) {
	msg := message.NewMessage(&message.MessageText{
		MessagerKernel: message.MessagerKernel{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			Type:             "text",
			To:               num,
		},
		Text: message.Text{
			PreviewUrl: false,
			Body:       "",
		},
	})

	resp := cl.SendMessage(msg)
	if resp.GetResponseError() == nil {
		msgID = resp.GetResponseSuccess().GetMessageId()
		return msgID
	}

	return msgID
}

func workerIsOnWhats(cl *ClientWA, wg *sync.WaitGroup, numberIn <-chan string, dOut chan<- InfoContact) {
	defer wg.Done()

	for number := range numberIn {
		id := getMsgPing(cl, number)
		if id != "" {
			ch := registerContact(id, number)
			// wait for channel
			for {
				select {
				case <-time.After(5 * time.Second):
					dOut <- InfoContact{
						ContactPhone: number,
						RecipientID:  number,
						IsOnWhats:    false,
						IsError:      true,
						Error:        fmt.Errorf("Error in workerIsOnWhats. Error is: timeout"),
						MsgError:     "Error in workerIsOnWhats. Error is: timeout",
					}
					return
				case receive := <-ch:
					dOut <- InfoContact{
						ContactPhone: receive.ContactPhone,
						RecipientID:  receive.RecipientID,
						IsOnWhats:    receive.IsOnWhats,
						IsError:      receive.IsError,
						Error:        receive.Error,
					}
					return
				}
			}
		} else {
			dOut <- InfoContact{
				ContactPhone: number,
				RecipientID:  number,
				IsOnWhats:    false,
				IsError:      true,
				Error:        fmt.Errorf("Error in workerIsOnWhats. msgID is empty"),
			}
			return
		}

	}
}

func (c *ClientWA) IsOnWhats(contactPhone []string) []InfoContact {
	cant := len(contactPhone)
	infoContacts = make(map[string]pair, 0)
	var wg sync.WaitGroup

	// create channel
	dIn := make(chan string, (cant/2)+1)
	dOut := make(chan InfoContact, (cant/2)+1)

	// create workers
	for range cant {
		wg.Add(1)
		go workerIsOnWhats(c, &wg, dIn, dOut)
	}

	// Send data to workers
	go func() {
		// close channel
		defer close(dIn)

		for _, phone := range contactPhone {
			dIn <- phone
		}
	}()

	// close channel in
	go func() {
		wg.Wait()

		close(dOut)
	}()

	// receive info of Channel
	infoContactsExit := make([]InfoContact, 0)
	for data := range dOut {
		infoContactsExit = append(infoContactsExit, data)
	}

	return infoContactsExit
}

// GetAllTemplate associated to a waba
func (c *ClientWA) GetAllTplFromWaba(waba_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s/%s", waba_id, "message_templates"), c.Config, RequestWithVersion, nil)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetAllTemplate request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetAllTemplate request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetAllTemplate request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b).GetTemplateResponse()
}

// GetAllTemplate pre approval from library, lands is a list of land codes ej: ["AR", "CO", "mx", "fr"] not sensitive to case
func (c *ClientWA) getAllTplFromLibrary(q ...QueryData) response.Responser {
	if len(q) == 0 {
		q = append(q, QueryData{})
	}

	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", "message_template_library"),
		c.Config, RequestWithQueryVersion, q[0])
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetAllTplFromLibrary request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetAllTplFromLibrary request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetAllTplFromLibrary request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b).GetTemplateResponse()
}

// GetAllTemplate pre approval from library, lands is a list of land codes ej: ["AR", "CO", "mx", "fr"] not sensitive to case
func (c *ClientWA) GetAllTplFromLibrary(q ...QueryData) response.Responser {
	res := c.getAllTplFromLibrary(q...)

	if res.GetResponseError() == nil {
		t := res.GetTemplateResponse()
		pag := t.Paging
		afterValue := pag.Cursors.After

		for {
			q[0].SetValue("after", afterValue)

			if pag.Next == "" {
				break
			}

			temp := c.getAllTplFromLibrary(q...)
			if v, ok := temp.(*response.TemplateResponse); v == nil && ok {
				break
			}

			if temp.GetResponseError() != nil {
				break
			}

			t_ := temp.GetTemplateResponse()
			pag = t_.Paging
			afterValue = pag.Cursors.After
			t.Data = append(t.Data, t_.Data...)
		}

		return t
	}

	return res
}

// GetTemplate by id
func (c *ClientWA) GetTplById(id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", id), c.Config, RequestWithVersion, nil)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetTemplateById request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetTemplateById request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetTemplateById request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b)
}

// GetTemplate by name
func (c *ClientWA) GetTplByName(name string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s/%s", c.getWabaID(), "message_templates"), c.Config, RequestWithQueryVersion, QueryData{"name": name})
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetTemplateName request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetTemplateName request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetTemplateName request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b)
}

func (c *ClientWA) SendReadNotification(messageID string) response.Responser {
	data := map[string]any{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageID,
	}

	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%s", "messages"), c.Config, data)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in SendReadNotification request of ClientWA. error is: ", err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in SendReadNotification request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "SendReadNotification", "ClientWA")
}

func (c *ClientWA) createUpdateTemplate(tpl *types.MockupTemplate, isUpdate bool) response.Responser {
	data := make(map[string]any)

	funcName := "CreateTemplate"
	if isUpdate {
		funcName = "UpdateTemplate"
	}

	if tpl == nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Error type: (%s) in %s function of ClientWA. error is: %s", types.MsgErrorUnrecognized, funcName, "Template is nil"),
		})
	}

	byt, err_ := json.Marshal(tpl)
	if err_ != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorParsingJson,
			Message: fmt.Sprintf("Error type: (%s) in %s function of ClientWA. error is: %s", types.MsgErrorParsingJson, funcName, err_.Error()),
		})
	}

	json.Unmarshal(byt, &data)
	if isUpdate {
		delete(data, "id")
		delete(data, "sub_category")
		delete(data, "parameter_format")
		delete(data, "status")
	}

	if c.getWabaID() == "" {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Error type: (%s) in %s function of ClientWA. error is: %s", types.MsgErrorUnrecognized, funcName, "WabaID is empty"),
		})
	}

	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%s", "message_templates"), c.Config, RequestWithWaba, data)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorInRequest,
			Message: fmt.Sprintf("Error type: (%s) in function %s of ClientWA. error is: %s", types.MsgErrorInRequest, funcName, err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorInRequestMeta,
			Message: fmt.Sprintf("Error type: (%s) in function %s of ClientWA. error is: %s", types.MsgErrorInRequestMeta, funcName, err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, fmt.Sprintf("%s", funcName), "ClientWA")
}

func (c *ClientWA) CreateTemplate(tpl *types.MockupTemplate) response.Responser {
	return c.createUpdateTemplate(tpl, false)
}

// UpdateTemplate updates a template by ID.
// If the Waba ID is empty, it returns an error.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) UpdateTemplate(tpl *types.MockupTemplate) response.Responser {
	return c.createUpdateTemplate(tpl, true)
}

type ParamDelete struct {
	ID, Name string
}

// DeleteTemplate deletes a template by ID or Name.
// If the Waba ID in client is empty, it returns an error.
// If the request is successful, the response is returned; otherwise, an error is returned.
// The ParamDelete struct should have the following fields:
// - ID: template id (Delete by id)
// - Name: template name (Delete by name)
// if both ID and Name are empty, it returns an error
func (c *ClientWA) DeleteTemplate(p ParamDelete) response.Responser {
	if c.getWabaID() == "" {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Error type: (%s) in CreateTemplate function of ClientWA. error is: %s", types.MsgErrorUnrecognized, "WabaID is empty"),
		})
	}

	if p.ID == "" && p.Name == "" {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Error type: (%s) in CreateTemplate function of ClientWA. error is: %s", types.MsgErrorUnrecognized, "ID and Name are empty"),
		})
	}

	data := map[string]any{
		"hsm_id": p.ID,
		"name":   p.Name,
	}

	if p.ID == "" {
		delete(data, "hsm_id")
	}

	if p.Name == "" {
		delete(data, "name")
	}

	_, _, err := defaultRequest(http.MethodDelete, fmt.Sprintf("/%s", "message_templates"), c.Config, RequestWithQueryWaba, (QueryData)(data))
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorInRequest,
			Message: fmt.Sprintf("Error type: (%s) in function CreateTemplate of ClientWA. error is: %s", types.MsgErrorInRequest, err.Error()),
		})
	}

	// Do request
	resp, err := doRequest(c.request, c)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorInRequestMeta,
			Message: fmt.Sprintf("Error type: (%s) in function CreateTemplate of ClientWA. error is: %s", types.MsgErrorInRequestMeta, err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "CreateTemplate", "ClientWA")
}
