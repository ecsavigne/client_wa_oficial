//lint:file-ignore ST1005 Ignore capitalized strings error
package wpp

import (
	"bytes"
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

	// "github.com/docker/docker-credential-helpers/client"

	"github.com/ecsavigne/client_wa_oficial/v2/client"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	evt_types "github.com/ecsavigne/client_wa_oficial/v2/types/general/response/event/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/message"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response/event"
	"golang.org/x/net/http2"
)

type (
	clientHttp struct {
		*http.Client
		BaseUrl *url.URL `json:"base_url"`
	}

	ClientWA struct {
		*Config    `json:"config"`
		typeClient client.CLIENT_TYPE
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
	// field := component.GetEntry()[0].GetChange()[0].Field
	field := component.GetNotificationType()
	phoneNumberID := ""

	if component.GetEntry()[0].GetChange()[0].GetValue().GetMetadata() != nil {
		phoneNumberID = component.GetEntry()[0].GetChange()[0].GetValue().GetMetadata().GetPhoneNumberID()
	}
	wabaID := component.GetEntry()[0].GetID()

	isForme = false
	if phoneNumberID == cl.Config.getphoneID() && field == "messages" {
		isForme = true
	} else if field == "messages" {
		isForme = false
	} else {
		if wabaID == cl.Config.getBusinessID() {
			isForme = true
		} else {
			if wabaID == cl.Config.getWabaID() {
				isForme = true
			}
		}
	}

	return isForme, evt_types.ParseTypeNotificationWebhook(field)
}

func (cl *ClientWA) Broadcast(data map[string]any) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	var evt event.EventInterface

	// Listener message of the server way WebHook
	message, err := json.Marshal(data)
	if err != nil {
		return
	}

	msg := codeWebHook(message)
	isForme, typeNotification := cl.messageIsForMe(msg)
	if !isForme {
		return
	}

	switch typeNotification {
	// Account
	case evt_types.WEBHOOK_NOTIFICATION_ACCOUNT_ALERTS:
		evt = &event.AccountAlertsEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_ACCOUNT_REVIEW_UPDATE:
		evt = &event.AccountReviewUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_ACCOUNT_UPDATE:
		evt = &event.AccountUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_BUSINESS_CAPABILITY_UPDATE:
		evt = &event.BusinessCapabilityUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_HISTORY:
		evt = &event.HistoryEvent{
			Components: msg,
		}

	case evt_types.WEBHOOK_NOTIFICATION_MESSAGE:
		// switch getTypeMessage(msg) {
		switch msg.GetTypeMessage() {
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
			// can be status message or another notification about message
			// status := getSatusMessage(msg)
			status := msg.GetSatusMessage()

			switch {
			case status != "":
				if isVailidStatusMessage(status) {
					evt = &event.StatusMessageEvent{
						Components: msg,
					}
				}
			default:
				cl.Config.EventHandle(message)
			}
		}
	// 	template
	case evt_types.WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_COMPONENTS_UPDATE:
		evt = &event.MessageTemplateComponentsUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_QUALITY_UPDATE:
		evt = &event.MessageTemplateQualityUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_STATUS_UPDATE:
		evt = &event.MessageTemplateStatusUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_PARTNER_SOLUTIONS:
		evt = &event.PartnerSolutionsEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_PAYMENT_CONFIGURATION_UPDATE:
		evt = &event.PaymentConfigurationUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_PHONE_NUMBER_NAME_UPDATE:
		evt = &event.PhoneNumberNameUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_PHONE_NUMBER_QUALITY_UPDATE:
		evt = &event.PhoneNumberQualityUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_SECURITY:
		evt = &event.SecurityEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_SMB_APP_STATE_SYNC:
		evt = &event.SmbAppStateSyncEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES:
		evt = &event.SmbMessageEchoesEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_TEMPLATE_CATEGORY_UPDATE:
		evt = &event.TemplateCategoryUpdateEvent{
			Components: msg,
		}
	case evt_types.WEBHOOK_NOTIFICATION_USER_PREFERENCES:
		evt = &event.UserPreferencesEvent{
			Components: msg,
		}
	default:
		cl.Config.EventHandle(message)
		return
	}

	cl.Config.EventHandle(evt)

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
		typeClient: client.CLIENT_WHATSAPP,
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

	// if c.WebhookSocket != "" && c.EventHandle != nil {
	// 	go cl.initWebHookSocket()
	// }

	return cl
}

func createClientHttp2(timeOut int) *http.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()

	tr.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
	http2.ConfigureTransport(tr)

	d := time.Duration(timeOut)
	return &http.Client{
		Timeout:   d * time.Second,
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

	c.pathPhone = path.Join(c.getVersion(), c.phoneID)

	c.pathWaba = path.Join(c.getVersion(), c.wabaID)

	c.pathBusiness = path.Join(c.getVersion(), c.businessID)

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
		if e := response.GetResponseRequest(res.Body, "doRequest", "400").GetResponseError(); e != nil {
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
		if e := response.GetResponseRequest(res.Body, "doRequest", "401").GetResponseError(); e != nil {
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
		if e := response.GetResponseRequest(res.Body, "doRequest", "404").GetResponseError(); e != nil {
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

	if !validTypeMsg(m, wpp.MessageTypeTemplate) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeTemplate, m.GetType()),
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
	if !validTypeMsg(m, wpp.MessageTypeText) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeText, m.GetType()),
		})
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("exception in SendText request of ClientWA. %s", e.Error()),
		})
	}

	return resp
}

// SendReaction sends a reaction message. It validates the message type to ensure it is a reaction.
// If the message type is incorrect, it returns an error. Otherwise, it makes a request to send the message.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) sendReaction(m message.Messager) response.Responser {
	if !validTypeMsg(m, wpp.MessageTypeReaction) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeReaction, m.GetType()),
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
	if !validTypeMsg(m, wpp.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(wpp.InteractiveTypeList) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", wpp.InteractiveTypeList),
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
	if !validTypeMsg(m.(*message.MessageInteractive), wpp.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(wpp.InteractiveTypeButtonResponse) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", wpp.InteractiveTypeButtonResponse),
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
	if !validTypeMsg(m.(*message.MessageInteractive), wpp.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(wpp.InteractiveTypeButtonUrl) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", wpp.InteractiveTypeButtonUrl),
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
	if !validTypeMsg(m.(*message.MessageInteractive), wpp.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(wpp.InteractiveTypeProcess) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", wpp.InteractiveTypeProcess),
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
	if !validTypeMsg(m.(*message.MessageInteractive), wpp.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(wpp.InteractiveTypeProduct) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", wpp.InteractiveTypeProduct),
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
	if !validTypeMsg(m.(*message.MessageInteractive), wpp.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(wpp.InteractiveTypeMultiProduct) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", wpp.InteractiveTypeMultiProduct),
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
	if !validTypeMsg(m.(*message.MessageInteractive), wpp.MessageTypeInteractive) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
		})
	} else {
		interactive, ok := m.(*message.MessageInteractive)
		if !ok {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
			})
		}
		if !interactive.IsType(wpp.InteractiveTypeCatalog) {
			return response.NewError(&response.Error{
				Type:    response.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("InteractiveProto.type must be '%s'", wpp.InteractiveTypeCatalog),
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
	if !validTypeMsg(m, wpp.MessageTypeAudio) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeAudio, m.GetType()),
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
	if !validTypeMsg(m, wpp.MessageTypeImage) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeImage, m.GetType()),
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
	if !validTypeMsg(m, wpp.MessageTypeVideo) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeVideo, m.GetType()),
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
	if !validTypeMsg(m, wpp.MessageTypeDocument) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeDocument, m.GetType()),
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
	if !validTypeMsg(m, wpp.MessageTypeSticker) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeSticker, m.GetType()),
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
	if !validTypeMsg(m, wpp.MessageTypeLocation) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeLocation, m.GetType()),
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
	if !validTypeMsg(m, wpp.MessageTypeContact) {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeContact, m.GetType()),
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

func (c *ClientWA) sendInteractive(m message.Messager) response.Responser {
	interactive := m.GetInteractiveMessage()
	if interactive != nil {
		switch interactive.GetInteractiveProto().Type {
		case wpp.InteractiveTypeList:
			return c.sendInteractiveList(m)
		case wpp.InteractiveTypeButtonResponse:
			return c.sendInteractiveButtonResponse(m)
		case wpp.InteractiveTypeProduct:
			c.sendInteractiveOneProduct(m)
		case wpp.InteractiveTypeMultiProduct:
			return c.sendInteractiveMultiProduct(m)
		case wpp.InteractiveTypeProcess:
			return c.sendInteractiveMsgProcess(m)
		case wpp.InteractiveTypeCatalog:
			return c.sendInteractiveCatalog(m)
		case wpp.InteractiveTypeButtonUrl:
			return c.sendInteractiveButtonUrl(m)
		}
	}
	return response.NewError(&response.Error{
		Type:    response.ResponseError,
		Code:    401,
		Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", wpp.MessageTypeInteractive, m.GetType()),
	})
}

func (c *ClientWA) SendMessage(m message.Messager) response.Responser {
	switch m.GetType() {
	case wpp.MessageTypeAudio:
		return c.sendAudioMessage(m)
	case wpp.MessageTypeContact:
		return c.sendContactMessage(m)
	case wpp.MessageTypeDocument:
		return c.sendDocumentMessage(m)
	case wpp.MessageTypeImage:
		return c.sendImageMessage(m)
	case wpp.MessageTypeInteractive:
		return c.sendInteractive(m)
	case wpp.MessageTypeLocation:
		return c.sendLocationMessage(m)
	case wpp.MessageTypeReaction:
		return c.sendReaction(m)
	case wpp.MessageTypeResponse:
		return c.sendResponseMsg(m)
	case wpp.MessageTypeSticker:
		return c.sendStickerMessage(m)
	case wpp.MessageTypeTemplate:
		return c.sendTemplate(m)
	case wpp.MessageTypeText:
		return c.sendTextMessage(m)
	case wpp.MessageTypeVideo:
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

	return response.JsonWrapperResponseRequest(b, response.ResponsePhone)
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

	// return response.JsonWrapperResponseRequest(b, response.ResponseWabaInfo)
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

// Register a number for use api
func (c *ClientWA) RegisterForUseApi(phone_id string) response.Responser {
	data := map[string]any{
		"messaging_product": "whatsapp",
		"pin":               "123456",
	}

	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%s/%s", phone_id, "register"), c.Config, RequestWithVersion, data)
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Error in RegisterForUseApi creating defaultRequest of ClientWA. error is: %s", err.Error()),
		})
	}

	// Do request
	byt, _ := io.ReadAll(c.request.Body)
	defer c.request.Body.Close()

	c.request.Body = io.NopCloser(bytes.NewReader(byt))

	resp, err := doRequest(c.request, c)
	if err != nil {

		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Error in RegisterForUseApi executing doRequest (%s {%s}) of ClientWA. error is: %s", c.request.URL.String(), string(byt), err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "RegisterForUseApi", "ClientWA", response.ResponseSuccess)
}

// Subscribed waba in apps for receive messages webhook
func (c *ClientWA) SubscribedWabaInApps(waba_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%s/%s", waba_id, "subscribed_apps"), c.Config, RequestWithVersion)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in SubscribedWabaInApps request of ClientWA. error is: ", err.Error()),
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
			Message: fmt.Sprintln("Error in SubscribedWabaInApps request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "SubscribedWabaInApps", "ClientWA", response.ResponseSuccess)
}

// GetInfoSubscribedWaba retrieves information about the applications that
// are subscribed to receive messages from a WhatsApp Business Account.
//   if subscribed waba
//   {
// 	"data": [
// 		{
// 			"whatsapp_business_api_data": {
// 				"link": "https://oficial.crmsocialhub.com.br/",
// 				"name": "socialhub",
// 				"id": "765398235585526"
// 			}
// 		}
// 	]
// }

// else subscribed waba

// 	{
// 		"data": []''
// 	}

// If the request fails, it returns a ResponserRequest with the error
func (c *ClientWA) GetInfoSubscribedWaba(waba_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s/%s", waba_id, "subscribed_apps"), c.Config, RequestWithVersion)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetInfoSubscribedWaba request of ClientWA. error is: ", err.Error()),
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
			Message: fmt.Sprintln("Error in GetInfoSubscribedWaba request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "GetInfoSubscribedWaba", "ClientWA", response.ResponseWABA)
}

func (c *ClientWA) UnSubscribedWaba(waba_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodDelete, fmt.Sprintf("/%s/%s", waba_id, "subscribed_apps"), c.Config, RequestWithVersion)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in UnSubscribedWaba request of ClientWA. error is: ", err.Error()),
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
			Message: fmt.Sprintln("Error in UnSubscribedWaba request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.GetResponseRequest(resp.Body, "UnSubscribedWaba", "ClientWA", response.ResponseSuccess)
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

	return response.JsonWrapperResponseRequest(b, response.ResponseWabaInfo)
}

// GetBusinessInfo retrieves detailed information about a business using its business ID.
// It performs a GET request to fetch fields such as id, name, extended updated time,
// link, two-factor type, is_hidden status, payment account ID, verification status,
// updated time, and created time. If an error occurs during the request, a response
// error is returned. The response is wrapped and returned as a ResponserRequest.

func (c *ClientWA) GetBusinessInfo(business_id string) response.Responser {
	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", business_id), c.Config, RequestWithQueryVersion, QueryData{"fields": "id,name,extended_updated_time,link,two_factor_type,is_hidden,payment_account_id,verification_status,updated_time,created_time,whatsapp_business_manager_messaging_limit"})
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

	return response.JsonWrapperResponseRequest(b, response.ResponseBusiness)
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

	// if k := wabas.GetResponseWaba(); k != nil {
	// 	if k.Data[0] == (response.WabaInfo{}) {
	// 		return nil, nil
	// 	}
	// }

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
	msg := message.NewMessageWpp(&message.MessageText{
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

func (c *ClientWA) getAllTplFromWaba(waba_id string, q ...QueryData) response.Responser {
	if len(q) == 0 {
		q = append(q, QueryData{})
	}

	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s/%s", waba_id, "message_templates"),
		c.Config, RequestWithQueryVersion, q[0])
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

	return response.JsonWrapperResponseRequest(b, response.ResponseTemplate).GetTemplateResponse()
}

// getAllTemplate associated to a waba
func (c *ClientWA) GetAllTplFromWaba(waba_id string) response.Responser {
	res := c.getAllTplFromWaba(waba_id)

	q := append([]QueryData{}, QueryData{})
	if res.GetResponseError() == nil {
		t := res.GetTemplateResponse()
		pag := t.Paging
		afterValue := pag.Cursors.After

		for {
			q[0].SetValue("after", afterValue)

			if pag.Next == "" {
				break
			}

			temp := c.getAllTplFromWaba(waba_id, q...)
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

	return response.JsonWrapperResponseRequest(b, response.ResponseTemplate).GetTemplateResponse()
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

	return response.JsonWrapperResponseRequest(b, response.ResponseTemplate)
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

func (c *ClientWA) createUpdateTemplate(tpl *wpp.MockupTemplate, isUpdate bool) response.Responser {
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

	return response.GetResponseRequest(resp.Body, fmt.Sprintf("%s", funcName), "ClientWA", response.ResponseTemplate)
}

func (c *ClientWA) CreateTemplate(tpl *wpp.MockupTemplate) response.Responser {
	return c.createUpdateTemplate(tpl, false)
}

// UpdateTemplate updates a template by ID.
// If the Waba ID is empty, it returns an error.
// If the request is successful, the response is returned; otherwise, an error is returned.
func (c *ClientWA) UpdateTemplate(tpl *wpp.MockupTemplate) response.Responser {
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

// Get permissions of a token.
func (c *ClientWA) DebugToken(token string) response.Responser {
	if token == "" {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DebugToken request of ClientWA. error is: Token is empty"),
		})
	}

	_, _, err := defaultRequest(http.MethodGet, fmt.Sprintf("/%s", "debug_token"), c.Config, RequestWithQueryVersion, QueryData{"input_token": token})
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DebugToken request of ClientWA. error is: ", err.Error()),
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
			Message: fmt.Sprintln("Error in DebugToken request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in DebugToken request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b, response.ResponseDebugToken)
}

// GetLimiteMsg makes a GET call to the WhatsApp API to get the
// limit of messages that can be sent from a WhatsApp Business
// Account. It returns a ResponserRequest with the response data in
// the Body field, which is a JSON object with the following structure:
//
//	{
//	    "whatsapp_business_manager_messaging_limit": TIER_2K,
//		"id": "129324348517543531"
//	}
//
// If the request fails, it returns a ResponserRequest with the error
func (c *ClientWA) GetLimiteMsg(business_id uint) response.Responser {
	if business_id == 0 {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetLimiteMsg request of ClientWA. error is: business_id is empty"),
		})
	}
	c.setBusinessID(fmt.Sprint(business_id))

	_, _, err := defaultRequest(http.MethodGet, "/", c.Config, RequestWithQueryBusiness, QueryData{"fields": "whatsapp_business_manager_messaging_limit"})
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetLimiteMsg request of ClientWA. error is: ", err.Error()),
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
			Message: fmt.Sprintln("Error in GetLimiteMsg request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in GetLimiteMsg request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b, response.ResponseOther)
}

// UnregisterNumber removes a phone number from the WhatsApp Business Account.
// It returns a ResponserRequest with the response data in the Body field, which is a JSON object with the following structure:
//
//	{
//	    "success": true,
//	}
//
// If the request fails, it returns a ResponserRequest with the error
func (c *ClientWA) UnregisterNumber(phone_id uint) response.Responser {
	if phone_id == 0 {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in UnregisterNumber request of ClientWA. error is: phone_id is empty"),
		})
	}
	c.setPhoneID(fmt.Sprint(phone_id))

	_, _, err := defaultRequest(http.MethodPost, fmt.Sprintf("/%d/%s", phone_id, "deregister"), c.Config, RequestWithVersion, nil)
	if err != nil {
		if err, ok := err.(*response.Error); ok {
			return err
		}
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in UnregisterNumber request of ClientWA. error is: ", err.Error()),
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
			Message: fmt.Sprintln("Error in UnregisterNumber request of ClientWA. error is: ", err.Error()),
		})
	}

	// prepare response
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response.NewError(&response.Error{
			Type:    response.ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintln("Error in UnregisterNumber request of ClientWA. error is: ", err.Error()),
		})
	}

	return response.JsonWrapperResponseRequest(b, response.ResponseSuccess)
}
