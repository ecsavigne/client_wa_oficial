//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/ecsavigne/client_wa_oficial/event"
	"github.com/ecsavigne/client_wa_oficial/types"

	"github.com/gorilla/websocket"
)

type clientHttp struct {
	*http.Client
	BaseUrl *url.URL `json:"base_url"`
}

type Config struct {
	Token string `json:"token"`
	// Path del archivo .env incluyendo el nombre del archivo sin la extensión ej: file: /.../../config_env.env -> EnvFilePath: /.../../config_env
	EnvFilePath string `json:"env_file_path"`
	Error       error
	// Url del servidor WebHook con ruta /ws para conectar con el servidor WebSocket ej: wss://webhooks.savcoe-services.com/ws
	WebhookSocket string    `json:"webhook_socket"`
	EventHandle   func(any) // Funcion para manejar los eventos del servidor WebHook WebSocket
	path          string
	clientHttp
	request *http.Request
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
		}
		cl.EventHandle(evt)
		return
		// cl.Error = fmt.Errorf("Error connecting to WebSocket error is : %v", err)
	}
	defer conn.Close()

	// Escuchar mensajes del servidor
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
	err := setEnv(c[0].EnvFilePath)
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

	c.path = path.Join(CLOUD_API_VERSION, WA_PHONE_NUMBER_ID)

	c.BaseUrl, _ = url.Parse(WA_BASE_URL)

	if c.Token == "" {
		c.Token = CLOUD_API_ACCESS_TOKEN
	}

	if c.Client == nil {
		c.Client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &c
}

func (c *ClientWA) makeRequest(methoth string, ePoint string, msg types.Messager) (*types.ResponseRequest, error) {
	if msg.GetMessageLink() != "" {
		multipartRequest(methoth, ePoint, c.Config, msg)
	} else {
		deafaultRequest(methoth, ePoint, c.Config, msg)
	}

	if c.Config.Error != nil {
		log := fmt.Sprintf("Error in function makeRequest creting request of ClientWA. error is: %s", c.Config.Error.Error())
		c.Config.Error = fmt.Errorf("%s", log)
		fmt.Println(log)
		return nil, c.Config.Error
	}

	res, e := c.clientHttp.Do(c.request)
	if e != nil {
		log := fmt.Sprintf("Error in function makeRequest executing request of ClientWA. Message type: %s. error is: %s", msg.GetType(), e.Error())
		c.Config.Error = fmt.Errorf("%s", log)
		fmt.Println(log)
		return nil, c.Config.Error
	}
	defer res.Body.Close()

	var responseReq *types.ResponseRequest
	bodyResponse, e := io.ReadAll(res.Body)
	if e != nil {
		log := fmt.Sprintf("Error in function makeRequest reading response of ClientWA. Message type: %s. error is: %s", msg.GetType(), e.Error())
		c.Config.Error = fmt.Errorf("%s", log)
		fmt.Println(log)
		return nil, c.Config.Error
	}

	err := json.Unmarshal(bodyResponse, &responseReq)
	if err != nil {
		log := fmt.Sprintf("Error in function makeRequest creating responseReq of ClientWA. Message type: %s. error is: %s", msg.GetType(), err.Error())
		fmt.Println(log)
		c.Config.Error = fmt.Errorf("%s", log)
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

func (c *ClientWA) SendTemplate(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeTemplate) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeTemplate, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendTemplate request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

// SendTextMessage send a text message
func (c *ClientWA) SendTextMessage(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeText) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeText, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendText request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendReaction(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeReaction) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeReaction, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendReaction request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendInteractiveList(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

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
		log = fmt.Sprintln("Error en SendInteractiveList request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendInteractiveButtonResponse(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

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
		log = fmt.Sprintln("Error en SendInteractiveButtonResponse request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendInteractiveButtonUrl(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

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
		log = fmt.Sprintln("Error en SendInteractiveButtonUrl request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendInteractiveMsgProcess(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

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
		log = fmt.Sprintln("Error en SendInteractiveMsgProcess request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendInteractiveOneProduct(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

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
		log = fmt.Sprintln("Error en SendInteractiveOneProduct request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendInteractiveMultiProduct(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

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
		log = fmt.Sprintln("Error en SendInteractiveMultiProduct request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendInteractiveCatalog(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

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
		log = fmt.Sprintln("Error en SendInteractiveCatalog request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendResponseMsg(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, "response") {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message expect '%s', but get '%s'", "response", m.GetType()),
		}
	} else {
		if m.(*types.MessageResponse).Type != types.MessageTypeText {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeText, m.(*types.MessageResponse).Type),
			}
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendResponseMsg request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendAudioMessage(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeAudio) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeAudio, m.GetType()),
		}
	} else {
		if m.(*types.MessageAudio).Link != "" {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: "Message.link can't be exist",
			}
		}
		if m.(*types.MessageAudio).Id == "" {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: "Message.id can't be empty",
			}
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendAudioMessage request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendImageMessage(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeImage) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeImage, m.GetType()),
		}
	} else {
		if m.(*types.MessageImage).Link != "" && m.(*types.MessageImage).Id != "" {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: "Expect Message.id or Message.link, but not both",
			}
		}
	}

	if m.(*types.MessageImage).Link != "" {
		resp, e := c.makeRequest(http.MethodPost, "/media", m)
		if e != nil {
			log = fmt.Sprintln("Error en SendImage request of ClientWA. error is: ", e.Error())
			fmt.Println(log)
			panic(log)
		} else if resp.Error != nil {
			return resp.Error
		}
		m.(*types.MessageImage).Id = resp.Id
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendImage request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendVideoMessage(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeVideo) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeVideo, m.GetType()),
		}
	} else {
		if m.(*types.MessageVideo).Link != "" && m.(*types.MessageVideo).Id != "" {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: "Expect Message.id or Message.link, but not both",
			}
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendVideoMessage request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendDocumentMessage(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeDocument) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeDocument, m.GetType()),
		}
	} else {
		if m.(*types.MessageDocument).Link != "" && m.(*types.MessageDocument).Id != "" {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: "Expect Message.id or Message.link, but not both",
			}
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendDocumentMessage request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendStickerMessage(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeSticker) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeSticker, m.GetType()),
		}
	} else {
		if m.(*types.MessageSticker).Link != "" {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: "Message.link can't be exist",
			}
		}
		if m.(*types.MessageSticker).Id == "" {
			return &types.Error{
				Type:    types.ResponseError,
				Code:    401,
				Message: "Message.id can't be empty",
			}
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendStickerMessage request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendLocationMessage(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeLocation) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeLocation, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendLocationMessage request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) SendContactMessage(m types.Messager) types.ResponserRequest {
	log := ""
	var r types.ResponserRequest

	if !validTypeMsg(m, types.MessageTypeContact) {
		return &types.Error{
			Type:    types.ResponseError,
			Code:    401,
			Message: fmt.Sprintf("Message.type expect '%s', but get '%s'", types.MessageTypeContact, m.GetType()),
		}
	}

	resp, e := c.makeRequest(http.MethodPost, "/messages", m)
	if e != nil {
		log = fmt.Sprintln("Error en SendContactMessage request of ClientWA. error is: ", e.Error())
		fmt.Println(log)
		panic(log)
	}

	if resp.Error != nil {
		r = resp.Error
	} else {
		r = resp.Success
	}

	return r
}

func (c *ClientWA) UploadFile(mt types.MediaType) {

}

func (c *ClientWA) DownloadFile() {

}
