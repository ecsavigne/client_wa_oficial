package ig

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ecsavigne/client_wa_oficial/v2/client"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	generalpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	gralpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	gralrequest "github.com/ecsavigne/client_wa_oficial/v2/types/general/request"
	"github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	gralresponse "github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	gralevt "github.com/ecsavigne/client_wa_oficial/v2/types/general/response/event"
	evt_types "github.com/ecsavigne/client_wa_oficial/v2/types/general/response/event/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/ig"
	igpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/ig/gen/igpb/v1"
	"github.com/ecsavigne/client_wa_oficial/v2/types/ig/response/event"
	"github.com/pusher/pusher-websocket-go"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type Config[C client.ClientReceiver, D client.Deserializer] struct {
	token  string
	userID string
	// Path del archivo .env incluyendo el nombre del archivo sin la extensión ej: file: /.../../config_env.env -> EnvFilePath: /.../../config_env
	envFilePath     string
	eventHandle     func(any) // Funcion para manejar los eventos del servidor WebHook WebSocket
	version         string    // v15.0
	baseUrlFacebook string    // https://graph.facebook.com
	baseUrlIG       string    // https://graph.instagram.com
	appId           string
	appSecret       string
	// allow receive messages from webhook server via pusher, centrifugo or kafka protocol
	clientReceiver C
	// Deserializer only for Kafka it can be nil,
	deserializer D
}

func (self *Config[C, D]) GetType() client.TYPE_CONFIG          { return client.TYPE_CONFIG_IG }
func (self *Config[C, D]) SetToken(token string)                { self.token = token }
func (self *Config[C, D]) GetToken() string                     { return self.token }
func (self *Config[C, D]) SetUserID(userID string)              { self.userID = userID }
func (self Config[C, D]) GetUserID() string                     { return self.userID }
func (self *Config[C, D]) SetEnvFilePath(envFilePath string)    { self.envFilePath = envFilePath }
func (self Config[C, D]) GetEnvFilePath() string                { return self.envFilePath }
func (self *Config[C, D]) SetEventHandle(eventHandle func(any)) { self.eventHandle = eventHandle }
func (self Config[C, D]) GetEventHandle() func(any)             { return self.eventHandle }
func (self *Config[C, D]) SetVersion(version string)            { self.version = version }
func (self *Config[C, D]) GetVersion() string                   { return self.version }
func (self *Config[C, D]) SetBaseUrlFacebook(baseUrlFacebook string) {
	self.baseUrlFacebook = baseUrlFacebook
}
func (self Config[C, D]) GetBaseUrlFacebook() string     { return self.baseUrlFacebook }
func (self *Config[C, D]) SetBaseUrlIG(baseUrlIG string) { self.baseUrlIG = baseUrlIG }
func (self Config[C, D]) GetBaseUrlIG() string           { return self.baseUrlIG }
func (self *Config[C, D]) SetBaseUrl(baseUrl string)     { self.baseUrlIG = baseUrl }
func (self *Config[C, D]) GetBaseUrl() string            { return self.baseUrlIG }
func (self *Config[C, D]) SetAppId(appID string)         { self.appId = appID }
func (self Config[C, D]) GetAppId() string               { return self.appId }
func (self *Config[C, D]) SetAppSecret(appSecret string) { self.appSecret = appSecret }
func (self Config[C, D]) GetAppSecret() string           { return self.appSecret }
func (self Config[C, D]) String() string {
	return fmt.Sprintf("Config{token: %s, userID: %s, envFilePath: %s, version: %s, baseUrlFacebook: %s, baseUrlIG: %s, appId: %s, appSecret: %s}",
		self.token, self.userID, self.envFilePath, self.version, self.baseUrlFacebook, self.baseUrlIG, self.appId, self.appSecret)
}

type OptionsClientIG[C client.ClientReceiver, D client.Deserializer] func(*Config[C, D]) // Opts

func WithToken[C client.ClientReceiver, D client.Deserializer](token string) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.token = token
	}
}

func WithUserID[C client.ClientReceiver, D client.Deserializer](userID string) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.userID = userID
	}
}

func WithEnvFilePath[C client.ClientReceiver, D client.Deserializer](envFilePath string) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.envFilePath = envFilePath
	}
}

func WithEventHandle[C client.ClientReceiver, D client.Deserializer](eventHandle func(any)) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.eventHandle = eventHandle
	}
}

func WithVersion[C client.ClientReceiver, D client.Deserializer](version string) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.version = version
	}
}

func WithBaseUrlFacebook[C client.ClientReceiver, D client.Deserializer](baseUrl string) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.baseUrlFacebook = baseUrl
	}
}

func WithBaseUrlIG[C client.ClientReceiver, D client.Deserializer](baseUrl string) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.baseUrlIG = baseUrl
	}
}

func WithAppId[C client.ClientReceiver, D client.Deserializer](appID string) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.appId = appID
	}
}

func WithAppSecret[C client.ClientReceiver, D client.Deserializer](appSecret string) OptionsClientIG[C, D] {
	return func(c *Config[C, D]) {
		c.appSecret = appSecret
	}
}

func setEnv[C client.ClientReceiver, D client.Deserializer](c *Config[C, D]) error {
	var envPath string = c.envFilePath
	pathDir := path.Dir(envPath)
	envName := path.Base(envPath)
	viper.AddConfigPath(pathDir)
	viper.SetConfigType("env")
	viper.SetConfigName(fmt.Sprintf("%s.env", envName))
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("\033[31mError: No encontrado archivo app.env ni .cobraToml de tipo (toml) en\033[30m %s\n", pathDir)
		return fmt.Errorf("code error: %d, type of error: %s, message: %s", types.CodeErrorEnvNotFound, types.TypeErrorConfig, types.MsgErrorEnvNotFound)
	} else {
		if v := viper.GetString("IG_API_VERSION"); v != "" {
			c.SetVersion(v)
		}

		if v := viper.GetString("IG_BASE_URL"); v != "" {
			c.SetBaseUrlIG(v)
		}

		if v := viper.GetString("M4D_APP_ID"); v != "" {
			c.SetAppId(v)
		}

		if v := viper.GetString("M4D_APP_SECRET"); v != "" {
			c.SetAppSecret(v)
		}

	}

	return nil
}

func defaultConfig[C client.ClientReceiver, D client.Deserializer]() *Config[C, D] {
	c := &Config[C, D]{
		token:           "",
		userID:          "",
		envFilePath:     "",
		eventHandle:     nil,
		version:         "v25.0",
		baseUrlFacebook: "https://graph.facebook.com",
		baseUrlIG:       "https://graph.instagram.com",
		appId:           "",
		appSecret:       "",
	}

	return c
}

type ClientIG[C client.ClientReceiver, D client.Deserializer] struct {
	config     *Config[C, D]
	typeClient client.CLIENT_TYPE
	url        *url.URL
}

func NewClientIG[C client.ClientReceiver, D client.Deserializer](opts ...OptionsClientIG[C, D]) (*ClientIG[C, D], error) {
	c := defaultConfig[C, D]()

	for _, opt := range opts {
		opt(c)
	}

	err := setEnv(c)
	if err != nil {
		return nil, err
	}

	cl := &ClientIG[C, D]{
		config:     c,
		typeClient: client.CLIENT_IG,
	}

	cl.listenWebHook()

	return cl, nil
}

func (self *ClientIG[C, D]) listenWebHook() {
	switch c := any(self.config.clientReceiver).(type) {
	case client.NotTyp:
		fmt.Println("Not implemte received message from webhook")
	case *centrifuge.Client:
		fmt.Println("implemte received message from webhook via centrifuge")
	case *pusher.Client:
		fmt.Println("implemte received message from webhook via pusher")
	case *kafka.Consumer:
		fmt.Println("implemte received message from webhook via kafka")
	default:
		_ = c
		fmt.Println("received not recognized message from webhook")
	}
}

func (self *ClientIG[C, D]) MessageIsForMe(webHookData []byte) (isForMe bool, typeNotification evt_types.TYPE_NOTIFICATION_WEBHOOK) {
	webHookMsg := new(igpbv1.InstagramWebhookEvent)
	_ = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(webHookData, webHookMsg)

	return true, evt_types.ParseTypeNotificationWebhook(webHookMsg.GetObject())
}

func (self *ClientIG[C, D]) GetTypeMessage(msg *igpbv1.InstagramWebhookEvent) (typ string) {
	defer func() {
		if r := recover(); r != nil {
			typ = ""
		}
	}()

	// return msg.GetEntry()[0].GetChanges()[0].GetValue().
	return "msg.Entry[0].Changes[0].Value.Messages[0].Type"
}

func (self *ClientIG[C, D]) IsVailidStatusMessage(status string) bool {
	if status == "read" || status == "delivered" || status == "sent" || status == "failed" ||
		status == "deleted" || status == "warning" {
		return true
	}

	return false
}

func (self *ClientIG[C, D]) GetSatusMessage(msg *igpbv1.InstagramWebhookEvent) (status string) {
	defer func() {
		if r := recover(); r != nil {
			status = ""
		}
	}()

	return " msg.Entry[0].Changes[0].Value.Statuses[0].Status"
}

func (self *ClientIG[C, D]) Broadcast(data map[string]any) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	var evt gralevt.EventInterface

	// Listener message of the server way WebHook
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	// msg := codeWebHook(message)
	msg := new(igpbv1.InstagramWebhookEvent)
	_ = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(dataBytes, msg)
	isForme, typeNotification := self.MessageIsForMe(dataBytes)
	if !isForme {
		return
	}

	switch typeNotification {
	// case len(msg.Entry) != 0 &&
	// 	len(msg.Entry[0].Changes) != 0 &&
	// 	len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
	// 	msg.Entry[0].Changes[0].Value.Messages[0].Type != ""
	// :
	case evt_types.WEBHOOK_NOTIFICATION_MESSAGE:
		switch self.GetTypeMessage(msg) {
		case "audio":
			evt = &event.IGMessageAudioEvent{
				InstagramWebhookEvent: msg,
			}
		case "button":
			evt = &event.IGMessageButtonEvent{
				InstagramWebhookEvent: msg,
			}
		case "document":
			evt = &event.IGMessageDocumentEvent{
				InstagramWebhookEvent: msg,
			}
		case "text":
			evt = &event.IGMessageTextEvent{
				InstagramWebhookEvent: msg,
			}
		case "image":
			evt = &event.IGMessageImageEvent{
				InstagramWebhookEvent: msg,
			}
		case "interactive":
			evt = &event.IGMessageInteractiveEvent{
				InstagramWebhookEvent: msg,
			}
		case "order":
			evt = &event.IGMessageOrderEvent{
				InstagramWebhookEvent: msg,
			}
		case "sticker":
			evt = &event.IGMessageStickerEvent{
				InstagramWebhookEvent: msg,
			}
		case "system":
			evt = &event.IGMessageSystemEvent{
				InstagramWebhookEvent: msg,
			}
		case "video":
			evt = &event.IGMessageVideoEvent{
				InstagramWebhookEvent: msg,
			}
		case "reaction":
			evt = &event.IGMessageReactionEvent{
				InstagramWebhookEvent: msg,
			}
		case "location":
			evt = &event.IGMessageLocationEvent{
				InstagramWebhookEvent: msg,
			}
		case "contacts":
			evt = &event.IGMessageContactEvent{
				InstagramWebhookEvent: msg,
			}
		case "unknown":
			evt = &event.IGMessageUnknownEvent{
				InstagramWebhookEvent: msg,
			}
		default:
			// can be status message or another notification about message
			status := self.GetSatusMessage(msg)

			switch {
			case status != "":
				if self.IsVailidStatusMessage(status) {
					evt = &event.IGStatusMessageEvent{
						InstagramWebhookEvent: msg,
					}
				}
			default:
				self.GetConfig().GetEventHandle()(msg)
			}
		}

	case evt_types.WEBHOOK_NOTIFICATION_TEMPLATE_CATEGORY_UPDATE,
		evt_types.WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_STATUS_UPDATE:
		evt = &event.IGMessageTemplateEvent{
			InstagramWebhookEvent: msg,
		}
	default:
		self.GetConfig().GetEventHandle()(msg)
		return
	}

	self.GetConfig().GetEventHandle()(evt)
}

func (self *ClientIG[C, D]) String() string {
	return fmt.Sprintf("ClientIG{config: %s, typeClient: %s, url_base: %s}", self.config.String(), self.typeClient.String(), self.GetConfig().GetBaseUrl())
}

func (self *ClientIG[C, D]) MultipartRequest(method string, data proto.Message, ePoint string) (*http.Request, error) {
	var (
		// config = self.GetConfig()
		e error
	)

	if !strings.HasPrefix(ePoint, "/") {
		e = fmt.Errorf("Error in multipartRequest, file: RequestConfig.go.Error is: EndPoint is not start with /")
		return nil, e
	}

	return gralrequest.UtilMultipartRequest(self, data, method, ePoint)
}

// func (self *ClientIG[C, D]) createUrl() {
// 	// return self.typeClient.String()
// }

func (self *ClientIG[C, D]) GetType() string                { return self.typeClient.String() }
func (self *ClientIG[C, D]) GetConfig() client.ConfigClient { return self.config }

func (self *ClientIG[C, D]) executeRequest(method string, ePoint string, data any, isID bool, resp_ ...gralresponse.ResponseType) gralresponse.Responser {
	respType := gralresponse.ResponseUnknow
	if len(resp_) > 0 {
		respType = resp_[0]
	}

	// if msg.GetMessageLink() != "" || msg.GetFileHeader() != nil {
	// 	multipartRequest(methoth, ePoint, c.Config, msg)
	// } else {
	req, err := gralrequest.DefaultRequest(self.GetConfig(), method, ePoint, data, isID)
	if err != nil {
		log := fmt.Sprintf("Error in function executeRequest creting request of ClientIG. error is: %s", err.Error())
		errorMsg := &gralpbv1.ResponseError{}
		errorMsg.SetType(types.TypeErrorInRequest)
		errorMsg.SetCode(types.CodeErrorInRequest)
		errorMsg.SetMessage(log)

		return gralresponse.NewResponse(errorMsg, gralresponse.ResponseError)
	}

	return gralrequest.Do(self, req, respType)
}

func (self *ClientIG[C, D]) SubscribeWebHook() response.Responser {
	data := map[string]any{
		"subscribed_fields": []string{"messages", "messaging_postbacks", "messaging_seen", "messaging_handover", "messaging_referral", "message_reactions", "standby", "comments", "live_comments", "mentions", "story_insights"},
	}

	return self.executeRequest(http.MethodPost, "/subscribed_apps", data, true, gralresponse.ResponseSuccess)
}

func (self *ClientIG[C, D]) UnsubscribeWebHook() gralresponse.Responser {
	return self.executeRequest(http.MethodDelete, "/subscribed_apps", nil, true, gralresponse.ResponseSuccess)
}

func (self *ClientIG[C, D]) sendTextMessage(msg *igpbv1.InstagramTextMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
}

func (self *ClientIG[C, D]) sendReactionMessage(msg *igpbv1.InstagramReactionMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
}

// upload media to Meta servers and get media_id if msg.GetFileHeader() exist, upload media to Meta servers and get media_id, then send media message with media_id
// if msg.GetFileHeader() == nil and exist msg.GetMessageLink() or msg.GetMessage().GetId() != "" then send media message with message link or media id
// if msg.GetFileHeader() == nil and msg.GetFileHeader() == nil and not exist msg.GetMessageLink() and msg.GetMessage().GetId() == "" then return error response with message "Media file is required for media message"() == "" and msg.GetMessage().GetId() == "" then return error response with message "Media file not found. Please provide a media file or a media link or a media id"
func (self *ClientIG[C, D]) sendMediaMessage(msg *igpbv1.InstagramMediaMessage) gralresponse.Responser {
	// if msg.GEtMessage().GetId()
	if msg.GetMessage().GetAttachment().GetPayload().GetUrl() == "" &&
		msg.GetMessage().GetAttachment().GetPayload().GetAttachmentId() == "" &&
		msg.GetFileDescriptor() == nil {
		errorMsg := &gralpbv1.ResponseError{}
		errorMsg.SetCode(types.CodeErrorInRequest)
		errorMsg.SetType(types.TypeErrorInRequest)
		errorMsg.SetMessage("Media file not found. Please provide a media file or a media link or a media id")
	}

	// if send url media message
	if msg.GetMessage().GetAttachment().GetPayload().GetUrl() != "" {
		msg.SetFileDescriptor(nil)
		msg.GetMessage().GetAttachment().GetPayload().SetAttachmentId("")
		return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
	}

	// if send media id media message
	if msg.GetMessage().GetAttachment().GetPayload().GetAttachmentId() != "" {
		msg.SetFileDescriptor(nil)
		msg.GetMessage().GetAttachment().GetPayload().SetUrl("")
		return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
	}

	if msg.GetMessage().GetAttachment().GetType() == ig.IG_ATTACHMENT_TYPE_LIKE_HEART.String() || msg.GetMessage().GetAttachments() != nil {
		return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
	}

	// if msg.GetFileHeader() != nil {

	return nil
}

// TODO: Implementar
func (self *ClientIG[C, D]) sendMediaShareMessage(msg *igpbv1.InstagramMediaShareMessage) gralresponse.Responser {
	return nil
}

/**
* @name: sendInstagramQuickRepliesMessage
* @description: Send Instagram Quick Replies message to Instagram
* @param {proto.Message} msg e.g: *igpbv1.InstagramQuickRepliesMessage
* @return gralresponse.Responser
 */
func (self *ClientIG[C, D]) sendInstagramQuickRepliesMessage(msg *igpbv1.InstagramQuickRepliesMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
}

func (self *ClientIG[C, D]) sendInstagramPersistentMenu(msg *igpbv1.InstagramPersistentMenuMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messenger_profile", msg, true)
}

func (self *ClientIG[C, D]) sendButtonGenericTemplateMessage(msg *igpbv1.InstagramTemplateButtonTemplate) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/me/messages", msg, false, gralresponse.SentMessageResponse)
}

/**
* @name: SendAction
* @description: Send action to instagram
* @param {string} scope_id the id of the user to send the action to
* @param {string} action the action to send. Can be "typing_on",  "typing_off" or "mark_seen"
* @return gralresponse.Responser
 */
func (self *ClientIG[C, D]) sendAction(scope_id, action string) gralresponse.Responser {
	if action != "typing_on" && action != "typing_off" && action != "mark_seen" {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Action not recognized. Action expect 'typing_on', 'typing_off' or 'mark_seen'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	msg := new(igpbv1.InstagramSenderAction)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)
	msg.SetRecipient(recipient)

	msg.SetSenderAction(action)

	return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)

	// return self.executeRequest(http.MethodGet, "/", types.QueryData{
	// 	"fields": "id,user_id,media_count,name,username,followers_count,follows_count,profile_picture_url",
	// }, gralresponse.ResponseInfoAccountBusiness)
}

func (self *ClientIG[C, D]) SendPresence(recipient_id, action string) gralresponse.Responser {
	return self.sendAction(recipient_id, action)
}

func (self *ClientIG[C, D]) MarkRead(recipient_id string) gralresponse.Responser {
	return self.sendAction(recipient_id, "mark_seen")
}

func (self *ClientIG[C, D]) sendInstagramIceBreakersMessage(msg *igpbv1.InstagramIceBreakersMessage) gralresponse.Responser {
	resp := self.getInfoAccountBusiness().GetResponse()
	infoAccount, ok := resp.(*igpbv1.InstagramInfoAccountBusinessResponse)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Error getting info account business. Response type is not InstagramInfoAccountBusinessResponse")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	return self.executeRequest(http.MethodPost, fmt.Sprintf("/%s/%s", infoAccount.GetId(), "messenger_profile"), msg, false)
}
func (self *ClientIG[C, D]) createInstagramWelcomeMessageFlows(msg *igpbv1.InstagramWelcomeMessageFlows) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, fmt.Sprintf("/%s/%s", "me", "welcome_message_flows"), msg, false)
}

func (self *ClientIG[C, D]) sendPrivateReplyMessage(msg *igpbv1.InstagramPrivateReplyMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
}

func (self *ClientIG[C, D]) sendInstagramHumanAgentMessage(msg *igpbv1.InstagramHumanAgentMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/me/messages", msg, false, gralresponse.SentMessageResponse)
}

/**
* @name: SendMessage
* @description: Send message to instagram
* @param {proto.Message} msg e.g: *igpbv1.InstagramTextMessage, *igpbv1.InstagramMediaMessage, *igpbv1.InstagramMediaShareMessage
* @return gralresponse.Responser
 */
func (self *ClientIG[C, D]) SendMessage(msg proto.Message) gralresponse.Responser {
	switch v := msg.(type) {
	case *igpbv1.InstagramTextMessage:
		return self.sendTextMessage(v)
	case *igpbv1.InstagramReactionMessage:
		return self.sendReactionMessage(v)
	case *igpbv1.InstagramMediaMessage:
		return self.sendMediaMessage(v)
	case *igpbv1.InstagramMediaShareMessage:
		return self.sendMediaShareMessage(v)
	case *igpbv1.InstagramQuickRepliesMessage:
		return self.sendInstagramQuickRepliesMessage(v)
	case *igpbv1.InstagramPersistentMenuMessage:
		return self.sendInstagramPersistentMenu(v)
	case *igpbv1.InstagramIceBreakersMessage:
		return self.sendInstagramIceBreakersMessage(v)
	case *igpbv1.InstagramWelcomeMessageFlows:
		return self.createInstagramWelcomeMessageFlows(v)
	case *igpbv1.InstagramPrivateReplyMessage:
		return self.sendPrivateReplyMessage(v)
	case *igpbv1.InstagramHumanAgentMessage:
		return self.sendInstagramHumanAgentMessage(v)
	case *igpbv1.InstagramTemplateButtonTemplate:
		switch v.GetMessage().GetAttachment().GetPayload().GetTemplateType() {
		case ig.IG_TEMPLATE_BUTTON.String(), ig.IG_TEMPLATE_GENERIC.String():
			return self.sendButtonGenericTemplateMessage(v)
		default:
			errorResponse := &gralpbv1.ResponseError{}
			errorResponse.SetCode(401)
			errorResponse.SetMessage("Message not recognized. Attachment type expect 'template_button' or 'generic'")
			return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
		}
		// return self.sendInstagramTemplateButtonTemplate(v)
	default:
		fmt.Printf("Error in SendMessage, file: IGClient.go. Message type not recognized. Type of message is: %T\n", v)
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Message not recognized. Message.type expect 'text', 'audio', 'image', 'video', 'document', 'sticker', 'location', 'contact', 'template', 'interactive', 'reaction'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}

// Creates
type resposeVerifyContainer struct {
	response gralresponse.Responser
}

func (self *ClientIG[C, D]) verifyContainer(idContainer string, out chan<- resposeVerifyContainer) {
	// return self.createImagePostContainer(msg)
	data := types.QueryData{
		"fields": "status_code",
	}
	resp_ := self.executeRequest(http.MethodGet, fmt.Sprintf("/%s", idContainer), data, false, gralresponse.InstagramFieldContainerResponse)
	containerResponse, ok := resp_.GetResponse().(*igpbv1.InstagramFieldContainerResponse)
	if !ok {
		out <- resposeVerifyContainer{response: resp_}
		return
	}

	for true {
		switch containerResponse.GetStatusCode() {
		case "FINISHED":
			out <- resposeVerifyContainer{response: resp_}
			return
		case "ERROR":
			errorResponse := &gralpbv1.ResponseError{}
			errorResponse.SetCode(401)
			errorResponse.SetMessage("Error in container. Status code is ERROR")
			respose := gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
			out <- resposeVerifyContainer{response: respose}
			return
		default:
			time.Sleep(20 * time.Second)
			resp_ = self.executeRequest(http.MethodGet, fmt.Sprintf("/%s", idContainer), data, false, gralresponse.InstagramFieldContainerResponse)
			containerResponse, ok = resp_.GetResponse().(*igpbv1.InstagramFieldContainerResponse)
		}
	}
}

func (self *ClientIG[C, D]) createContainer(imgs, videos []string, mediaType ig.IG_MEDIA_TYPE) string {
	msg := new(igpbv1.InstagramContainerMessage)

	if mediaType == ig.IG_MEDIA_TYPE_STORIES {
		msg.SetMediaType(mediaType.String())
	}

	if mediaType == ig.IG_MEDIA_TYPE_CAROUSEL {
		msg.SetIsCarouselItem(true)
		// msg.SetMediaType(mediaType.String())
	}

	if len(imgs) > 0 {
		msg.SetImageUrl(imgs[0])
	} else if len(videos) > 0 {
		if mediaType == ig.IG_MEDIA_TYPE_POST || mediaType == ig.IG_MEDIA_TYPE_CAROUSEL {
			msg.SetMediaType(ig.IG_MEDIA_TYPE_REELS.String())
		}

		msg.SetVideoUrl(videos[0])
		msg.SetUploadType("resumable")
	}

	resp := self.executeRequest(http.MethodPost, fmt.Sprintf("/%s", "media"), msg, true, gralresponse.InstagramFieldContainerResponse)
	contResp, ok := resp.GetResponse().(*igpbv1.InstagramFieldContainerResponse)
	if !ok {
		return ""
	}

	return contResp.GetId()
}

type dataWorkerPreparedContainer struct {
	imgs         []string
	videos       []string
	mediaType    ig.IG_MEDIA_TYPE
	id_container string
}

func (self *ClientIG[C, D]) workerPreparedContainer(in <-chan dataWorkerPreparedContainer, out chan<- dataWorkerPreparedContainer) {
	for el := range in {
		fmt.Println("<<<<<< workerPreparedContainer >>>>>>  imgs: ", el.imgs, " videos: ", el.videos, " mediaType: ", el.mediaType)
		id_container := self.createContainer(el.imgs, el.videos, el.mediaType)
		out <- dataWorkerPreparedContainer{id_container: id_container}
	}
}

func (self *ClientIG[C, D]) workerVerifyContainers(in <-chan dataWorkerPreparedContainer, out chan<- dataWorkerPreparedContainer) {
	for el := range in {
		ch := make(chan resposeVerifyContainer, 1)
		go self.verifyContainer(el.id_container, ch)
		resp := <-ch
		if resp.response.GetType() == gralresponse.ResponseError {
			out <- dataWorkerPreparedContainer{id_container: ""}
		} else {
			out <- dataWorkerPreparedContainer{id_container: el.id_container}
		}
	}

}

func (self *ClientIG[C, D]) verifyContainers(cointainers_id []string) []string {
	var (
		cantEl                = len(cointainers_id)
		in                    = make(chan dataWorkerPreparedContainer, cantEl)
		out                   = make(chan dataWorkerPreparedContainer, cantEl)
		wg                    = new(sync.WaitGroup)
		container_id_verified = make([]string, 0, cantEl)
	)
	workers := math.Ceil(math.Sqrt(float64(cantEl)))
	workers = max(workers, 2)

	for range int(workers) {
		wg.Go(func() {
			// Verificar container
			self.workerVerifyContainers(in, out)
		})
	}

	//  send data to worker
	go func() {
		defer close(in)
		for _, v := range cointainers_id {
			in <- dataWorkerPreparedContainer{id_container: v}
		}
	}()

	go func() {
		defer close(out)
		wg.Wait()
	}()

	for el := range out {
		// get container verified
		container_id_verified = append(container_id_verified, el.id_container)
	}

	return container_id_verified
}

func (self *ClientIG[C, D]) createCaruselContainer(ids []string, caption string) (container_id string) {
	data := new(igpbv1.InstagramContainerMessage)
	data.SetMediaType("CAROUSEL")
	data.SetChildren(ids)
	if caption != "" {
		data.SetCaption(caption)
	}

	resp := self.executeRequest(http.MethodPost, "/media", data, true, gralresponse.InstagramFieldContainerResponse)

	contResp, ok := resp.GetResponse().(*igpbv1.InstagramFieldContainerResponse)
	if !ok {
		return ""
	}

	return contResp.GetId()
}

func (self *ClientIG[C, D]) preparedCaruselContainer(imgs, videos []string, mediaType ig.IG_MEDIA_TYPE, caption string) (container_id string) {
	var (
		cantEl = len(imgs) + len(videos)
		in     = make(chan dataWorkerPreparedContainer, cantEl)
		out    = make(chan dataWorkerPreparedContainer, cantEl)
		wg     = new(sync.WaitGroup)
	)
	workers := math.Ceil(math.Sqrt(float64(cantEl)))
	workers = max(workers, 2)

	for range int(workers) {
		wg.Go(func() {
			self.workerPreparedContainer(in, out)
		})
	}

	//  send data to worker
	go func() {
		defer close(in)

		for _, v := range imgs {
			in <- dataWorkerPreparedContainer{imgs: []string{v}, videos: []string{}, mediaType: mediaType}
		}
		for _, v := range videos {
			in <- dataWorkerPreparedContainer{imgs: []string{}, videos: []string{v}, mediaType: mediaType}
		}
	}()

	go func() {
		defer close(out)
		wg.Wait()
	}()

	arrContainerID := make([]string, 0, cantEl)
	for el := range out {
		arrContainerID = append(arrContainerID, el.id_container)
	}

	// verificar elementos do carrusel container e criar carrusel container
	arrContainerIDTemp := self.verifyContainers(arrContainerID)

	return self.createCaruselContainer(arrContainerIDTemp, caption)
}

func (self *ClientIG[C, D]) createInstagramPostHistory(dataParam map[string]any, mediaType ig.IG_MEDIA_TYPE) gralresponse.Responser {
	if len(dataParam) == 0 {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to create container. Please provide data in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	imgs, ok := dataParam["images_url"].([]string)
	videos, ok2 := dataParam["videos_url"].([]string)
	id_container := ""

	var (
		lenImgs, lenVideos int = 0, 0
		out                    = make(chan resposeVerifyContainer, 1)
	)

	if ok {
		lenImgs = len(imgs)
	}

	if ok2 {
		lenVideos = len(videos)
	}

	cant := lenImgs + lenVideos

	if (!ok && !ok2) || cant == 0 {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to create container. Please provide []images_url or []videos_url array string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	if cant > 10 {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("The maximum number of media that can be included in a post is 10. Please provide a maximum of 10 media in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	if cant > 1 && cant <= 10 {
		// criar carrucel container and get container id carrucel
		caption := ""
		val, ok := dataParam["caption"]
		if ok {
			caption, ok = val.(string)
			if !ok {
				caption = ""
			}
		}

		if mediaType == ig.IG_MEDIA_TYPE_POST {
			mediaType = ig.IG_MEDIA_TYPE_CAROUSEL
		} else {
			mediaType = ig.IG_MEDIA_TYPE_STORIES
		}

		id_container = self.preparedCaruselContainer(imgs, videos, mediaType, caption)
		if id_container == "" {
			errorResponse := &gralpbv1.ResponseError{}
			errorResponse.SetCode(401)
			errorResponse.SetMessage("Error creating carousel container. Please check the data provided and try again")
			return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
		}
	}

	if cant == 1 {
		// criar container and get container id
		// media_type: "" <=> POST, STORIES, VIDEO, REEL
		id_container = self.createContainer(imgs, videos, mediaType)
		if id_container == "" {
			errorResponse := &gralpbv1.ResponseError{}
			errorResponse.SetCode(401)
			errorResponse.SetMessage("Error creating container. Please check the data provided and try again")
			return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
		}
	}

	// verificar status do container até que seja FINISHED ou ERROR
	go self.verifyContainer(id_container, out)
	verify := <-out
	if verify.response.GetType() != gralresponse.InstagramFieldContainerResponse {
		return verify.response
	}

	data := map[string]any{"creation_id": id_container}
	return self.executeRequest(http.MethodPost, "/media_publish", data, true, gralresponse.InstagramFieldContainerResponse)
}

/*
@ig_media_id: id do post, stories or midia

	@data in Create: map[string]any{
		"message": "string", // comment message
		"ig_media_id": "string", // id do post, stories or midia
	}
*/
func (self *ClientIG[C, D]) createComment(data map[string]any) gralresponse.Responser {
	if data["message"] == nil || data["ig_media_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to create comment. Please provide 'message' and 'ig_media_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	msg, ok := data["message"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to create comment. 'message' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igMediaID, ok := data["ig_media_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to create comment. 'ig_media_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	queryData := types.QueryData{
		"message": msg,
	}

	return self.executeRequest(http.MethodPost, fmt.Sprintf("/%s/comments?%s", igMediaID, queryData.String()), nil, false, gralresponse.InstagramFieldContainerResponse)
}

/*
@ig_comment_id: id do comment

	@data in Create: map[string]any{
		"message": "string", // comment message
		"ig_comment_id": "string", // id do comment
	}
*/
func (self *ClientIG[C, D]) replyComment(data map[string]any) gralresponse.Responser {
	if data["message"] == nil || data["ig_comment_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to create comment. Please provide 'message' and 'ig_comment_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	msg, ok := data["message"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to create comment. 'message' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igCommentID, ok := data["ig_comment_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to create comment. 'ig_comment_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	queryData := types.QueryData{
		"message": msg,
	}

	return self.executeRequest(http.MethodPost, fmt.Sprintf("/%s/replies?%s", igCommentID, queryData.String()), nil, false, gralresponse.InstagramFieldContainerResponse)
}

/*
@ig_comment_id: id do comment

	@data in Create: map[string]any{
		"ig_comment_id": "string", // id do comment, comment is nivel 1 or comment reply
		"hide": bool, // true to hide comment, false to show comment
	}
*/
func (self *ClientIG[C, D]) hideShowComment(data map[string]any) gralresponse.Responser {
	if data["hide"] == nil || data["ig_comment_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to hide/show comment. Please provide 'hide' and 'ig_comment_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	hide, ok := data["hide"].(bool)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to hide/show comment. 'hide' must be a boolean in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igCommentID, ok := data["ig_comment_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to hide/show comment. 'ig_comment_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	queryData := map[string]any{
		"hide": hide,
	}

	return self.executeRequest(http.MethodPost, fmt.Sprintf("/%s", igCommentID), queryData, false, gralresponse.ResponseSuccess)
}

/*
@ig_media_id: id do post, stories or midia

	@data in Create: map[string]any{
		"ig_media_id": "string", // id do post, stories or midia
		"comment_enabled": bool, // true to enable comment, false to disable comment
	}
*/
func (self *ClientIG[C, D]) enableComment(data map[string]any) gralresponse.Responser {
	if data["comment_enabled"] == nil || data["ig_media_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to enable comment. Please provide 'comment_enabled' and 'ig_media_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	commentEnabled, ok := data["comment_enabled"].(bool)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to enable comment. 'comment_enabled' must be a boolean in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igMediaID, ok := data["ig_media_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to enable comment. 'ig_media_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	queryData := map[string]any{
		"comment_enabled": commentEnabled,
	}

	return self.executeRequest(http.MethodPost, fmt.Sprintf("/%s", igMediaID), queryData, false, gralresponse.ResponseSuccess)
}

func (self *ClientIG[C, D]) Create(typeCreate ig.IG_CREATE_TYPE, data ...map[string]any) gralresponse.Responser {
	dataParam := make(map[string]any)
	if len(data) > 0 {
		dataParam = data[0]
	}

	switch typeCreate {
	case ig.IG_CREATE_POST:
		return self.createInstagramPostHistory(dataParam, ig.IG_MEDIA_TYPE_POST)
	case ig.IG_CREATE_STORY:
		return self.createInstagramPostHistory(dataParam, ig.IG_MEDIA_TYPE_STORIES)
	case ig.IG_CREATE_COMMENT:
		return self.createComment(dataParam)
	case ig.IG_CREATE_REPLY_COMMENT:
		return self.replyComment(dataParam)
	case ig.IG_CREATE_HIDE_COMMENT:
		return self.hideShowComment(dataParam)
	case ig.IG_CREATE_ENABLE_COMMENT:
		return self.enableComment(dataParam)
	default:
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("typeCreate not recognized. typeCreate expect 'post', 'story', 'comment', 'reply_comment', 'hide_comment' or 'enable_comment'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}

// Gets
func (self *ClientIG[C, D]) getInfoAccountBusiness() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/", types.QueryData{
		"fields": "id,user_id,media_count,account_type,name,username,followers_count,follows_count,profile_picture_url",
	}, true, gralresponse.InfoAccountBusinessResponse)
}

func (self *ClientIG[C, D]) getInstagramPersistentMenu() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/messenger_profile", types.QueryData{
		"fields": "persistent_menu",
	}, true)
}

func (self *ClientIG[C, D]) getInstagramIceBreakers() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/messenger_profile", types.QueryData{
		"fields": "ice_breakers",
	}, true)
}

func (self *ClientIG[C, D]) getInstagramLink() gralresponse.Responser {
	resp := self.getInfoAccountBusiness().GetResponse()
	infoAccount, ok := resp.(*igpbv1.InstagramInfoAccountBusinessResponse)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Error getting info account business. Response type is not InstagramInfoAccountBusinessResponse")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	link := fmt.Sprintf("https://ig.me/%s", infoAccount.GetUsername())

	unknow := &generalpbv1.UnknownResponse{}
	structValue, _ := structpb.NewValue(map[string]any{"link": link})
	unknow.SetData(structValue)

	return gralresponse.NewResponse(unknow, gralresponse.ResponseUnknow)
}

func (self *ClientIG[C, D]) getWelcomeMessageFlowsADS() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/me/welcome_message_flows", types.QueryData{}, false)
}

/*
@ig_media_id: id do post, stories or midia

	@data in Get: map[string]any{
		"ig_media_id": "string", // id do post, stories ou midia
	}

@return: list of replies comments of the comment with from,text,timestamp,user,username,replies,parent_id,like_count,hidden
*/
func (self *ClientIG[C, D]) getComments(data map[string]any) gralresponse.Responser {
	if data["ig_media_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get comments. Please provide 'ig_media_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igMediaID, ok := data["ig_media_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get comments. 'ig_media_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	dataQuery := types.QueryData{
		"fields": "from,text,timestamp,user,username,replies,parent_id,like_count,hidden",
	}
	return self.executeRequest(http.MethodGet, fmt.Sprintf("/%s/comments", igMediaID), dataQuery, false, gralresponse.InstagramCommentResponse)
}

/*
@ig_comment_id: id do comment

	@data in Get: map[string]any{
		"ig_comment_id": "string", // id do comment
	}

@return: list of replies comments of the comment with text,timestamp,from,user,username
*/
func (self *ClientIG[C, D]) getRepliesComments(data map[string]any) gralresponse.Responser {
	if data["ig_comment_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get replies comments. Please provide 'ig_comment_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igCommentID, ok := data["ig_comment_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get replies comments. 'ig_comment_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	dataQuery := types.QueryData{
		"fields": "text,timestamp,from,user,username",
	}
	return self.executeRequest(http.MethodGet, fmt.Sprintf("/%s/replies", igCommentID), dataQuery, false, gralresponse.InstagramCommentResponse)
}

func (self *ClientIG[C, D]) getSubscibeWebhookField() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/subscribed_apps", nil, true, gralresponse.ResponseUnknow)
}

/*
IG_GET_METRICS_MIDIA              IG_GET_INFO_TYPE = "metrics_midia"

	IG_GET_METRICS_MIDIA_INSIGHT      IG_GET_INFO_TYPE = "metrics_midia_insight"
	IG_GET_METRICS_USER_INSIGHT
*/
func (self *ClientIG[C, D]) getMetricsMedia(data map[string]any) gralresponse.Responser {
	if data["ig_media_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get metrics media. Please provide 'ig_media_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igMediaID, ok := data["ig_media_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get metrics media. 'ig_media_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	dataQuery := types.QueryData{
		"fields": "like_count,comments_count,id,media_type",
	}
	return self.executeRequest(http.MethodGet, fmt.Sprintf("/%s", igMediaID), dataQuery, false, gralresponse.InstagramMetricResponse)
}

func (self *ClientIG[C, D]) getMetricsUserInsight() gralresponse.Responser {
	dataQuery := types.QueryData{
		"metric": "reach,follower_count,website_clicks,profile_views,online_followers,accounts_engaged,total_interactions,likes,comments,shares,saves,replies,engaged_audience_demographics,reached_audience_demographics,follower_demographics,follows_and_unfollows,profile_links_taps,views,threads_likes,threads_replies,reposts,quotes,threads_followers,threads_follower_demographics,content_views,threads_views,threads_clicks,threads_reposts",
		"period": "day",
	}

	return self.executeRequest(http.MethodGet, "/insights", dataQuery, true, gralresponse.InstagramMetricInsightResponse)
}

func (self *ClientIG[C, D]) getMetricsMediaInsight(data map[string]any) gralresponse.Responser {
	if data["ig_media_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get metrics media insight. Please provide 'ig_media_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igMediaID, ok := data["ig_media_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get metrics media insight. 'ig_media_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	dataQuery := types.QueryData{
		"metric": "shares,comments,likes,saved,total_interactions,follows,profile_visits,profile_activity,reach,views,content_views",
	}

	return self.executeRequest(http.MethodGet, fmt.Sprintf("/%s/insights", igMediaID), dataQuery, false, gralresponse.InstagramMetricInsightResponse)
}

func (self *ClientIG[C, D]) getListConversation() gralresponse.Responser {
	dataQuery := types.QueryData{
		"platform": "INSTAGRAM",
	}

	return self.executeRequest(http.MethodGet, "/me/conversations", dataQuery, false, gralresponse.InstagramListConversationResponse)
}

func (self *ClientIG[C, D]) getConversationWithUser(data map[string]any) gralresponse.Responser {
	if data["ig_user_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get conversation with user. Please provide 'ig_user_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	ig_user_id, ok := data["ig_user_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get conversation with user. Please provide 'ig_user_id' in data parameter and it must be a string")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	dataQuery := types.QueryData{
		"user_id": ig_user_id,
	}

	return self.executeRequest(http.MethodGet, "/me/conversations", dataQuery, false, gralresponse.InstagramListConversationResponse)
}

func (self *ClientIG[C, D]) getMessagesOfConversation(data map[string]any) gralresponse.Responser {
	if data["conversation_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get messages of conversation. Please provide 'conversation_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	conversation_id, ok := data["conversation_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get messages of conversation. Please provide 'conversation_id' in data parameter and it must be a string")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	dataQuery := types.QueryData{
		"fields": "messages{from,to,created_time,updated_at,message}",
	}

	return self.executeRequest(http.MethodGet, fmt.Sprintf("/%s", conversation_id), dataQuery, false, gralresponse.InstagramConversationMessageResponse)
}

func (self *ClientIG[C, D]) getInfoAboutMessage(data map[string]any) gralresponse.Responser {
	if data["message_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get info about message. Please provide 'message_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	message_id, ok := data["message_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to get info about message. Please provide 'message_id' in data parameter and it must be a string")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	dataQuery := types.QueryData{
		"fields": "id,created_time,from,to,message",
	}

	return self.executeRequest(http.MethodGet, fmt.Sprintf("/%s", message_id), dataQuery, false, gralresponse.ConversationMessageResponse)
}

func (self *ClientIG[C, D]) Get(type_info ig.IG_GET_INFO_TYPE, data ...map[string]any) gralresponse.Responser {
	dataParam := make(map[string]any)
	if len(data) > 0 {
		dataParam = data[0]
	}

	switch type_info {
	case ig.IG_GET_INFO_ACCOUNT_BUSINESS:
		return self.getInfoAccountBusiness()
	case ig.IG_GET_INFO_PERSISTENT_MENU:
		return self.getInstagramPersistentMenu()
	case ig.IG_GET_INFO_ICE_BREAKERS:
		return self.getInstagramIceBreakers()
	case ig.IG_GET_INFO_LINK:
		return self.getInstagramLink()
	case ig.IG_GET_INFO_WELCOME_MESSAGE_FLOWS:
		return self.getWelcomeMessageFlowsADS()
	case ig.IG_GET_COMMENT:
		return self.getComments(dataParam)
	case ig.IG_GET_REPLIES_COMMENTS:
		return self.getRepliesComments(dataParam)
	case ig.IG_GET_SUBSCRIBE_WEBHOOK_FIELD:
		return self.getSubscibeWebhookField()
	case ig.IG_GET_METRICS_MEDIA:
		return self.getMetricsMedia(dataParam)
	case ig.IG_GET_METRICS_MEDIA_INSIGHT:
		return self.getMetricsMediaInsight(dataParam)
	case ig.IG_GET_METRICS_USER_INSIGHT:
		return self.getMetricsUserInsight()
	case ig.IG_GET_LIST_CONVERSATION:
		return self.getListConversation()
	case ig.IG_GET_USER_CONVERSATION:
		return self.getConversationWithUser(dataParam)
	case ig.IG_GET_MESSAGES_CONVERSATION:
		return self.getMessagesOfConversation(dataParam)
	case ig.IG_GET_INFO_MESSAGE:
		return self.getInfoAboutMessage(dataParam)
	default:
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("type_info not recognized. type_info expect 'account_business', 'persistent_menu', 'ice_breakers', 'link', 'welcome_message_flows'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}

// Deletes
func (self *ClientIG[C, D]) deletePersistentMenu() gralresponse.Responser {
	return self.executeRequest(http.MethodDelete, "/messenger_profile", types.QueryData{
		"fields": []string{"persistent_menu"},
	}, true)
}

func (self *ClientIG[C, D]) deleteIceBreakers() gralresponse.Responser {
	return self.executeRequest(http.MethodDelete, "/messenger_profile", types.QueryData{
		"fields": []string{"ice_breakers"},
	}, true)
}

func (self *ClientIG[C, D]) deleteWelcomeMessageFlowsADS(data types.QueryData) gralresponse.Responser {
	if id, ok := data["flow_id"]; ok && id != "" {
		id = fmt.Sprintf("%c%s", '?', data.String())
		return self.executeRequest(http.MethodDelete, fmt.Sprintf("/me/welcome_message_flows%s", id), types.QueryData{}, false)
	}

	errorResponse := &gralpbv1.ResponseError{}
	errorResponse.SetCode(401)
	errorResponse.SetMessage("flow_id is required to delete welcome message flow. Please provide flow_id in data with key 'flow_id'")
	return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
}

/*
@ig_comment_id: id do comment

	@data in Delete: map[string]any{
		"ig_comment_id": "string", // id do comment to delete
	}
*/
func (self *ClientIG[C, D]) deleteComment(data map[string]any) gralresponse.Responser {
	if data["ig_comment_id"] == nil {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to delete comment. Please provide 'ig_comment_id' in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	igCommentID, ok := data["ig_comment_id"].(string)
	if !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to delete comment. 'ig_comment_id' must be a string in data parameter")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	return self.executeRequest(http.MethodDelete, fmt.Sprintf("/%s", igCommentID), nil, false, gralresponse.ResponseSuccess)
}

func (self *ClientIG[C, D]) Delete(typeDelete ig.IG_DELETE_TYPE, data ...map[string]any) gralresponse.Responser {
	dataParam := make(map[string]any)
	if len(data) > 0 {
		dataParam = data[0]
	}

	switch typeDelete {
	case ig.IG_DELETE_PERSISTENT_MENU:
		return self.deletePersistentMenu()
	case ig.IG_DELETE_ICE_BREAKERS:
		return self.deleteIceBreakers()
	case ig.IG_DELETE_WELCOME_MESSAGE_FLOWS:
		return self.deleteWelcomeMessageFlowsADS(dataParam)
	case ig.IG_DELETE_COMMENT:
		return self.deleteComment(dataParam)
	default:
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("typeDelete not recognized. typeDelete expect 'persistent_menu', 'ice_breakers', 'welcome_message_flows', 'comment'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}

// Updates
func (self *ClientIG[C, D]) updateWelcomeMessageFlowsADS(data types.QueryData) gralresponse.Responser {
	if len(data) == 0 {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Data is required to update welcome message flow. Please provide data with at least one of the following keys: 'flow_id', 'name', 'welcome_message', 'platforms'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	if _, ok := data["flow_id"]; !ok {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("flow_id is required to update welcome message flow. Please provide flow_id in data with key 'flow_id'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}

	if id, ok := data["flow_id"]; ok && id != "" {
		id = fmt.Sprintf("%cflow_id=%s", '?', id)
		var msg *igpbv1.InstagramWelcomeMessageFlows
		if msg_, ok := data["msg"]; ok {
			msg, ok = msg_.(*igpbv1.InstagramWelcomeMessageFlows)
			if !ok {
				errorResponse := &gralpbv1.ResponseError{}
				errorResponse.SetCode(401)
				errorResponse.SetMessage("msg is not of type InstagramWelcomeMessageFlows. Please provide msg in data with key 'msg' and of type InstagramWelcomeMessageFlows")
				return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
			}

			return self.executeRequest(http.MethodPost, fmt.Sprintf("/me/welcome_message_flows%s", id), msg, false)
		} else {
			errorResponse := &gralpbv1.ResponseError{}
			errorResponse.SetCode(401)
			errorResponse.SetMessage("msg is required to update welcome message flow. Please provide msg in data with key 'msg'")
			return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
		}
	}

	errorResponse := &gralpbv1.ResponseError{}
	errorResponse.SetCode(401)
	errorResponse.SetMessage("flow_id is required to update welcome message flow. Please provide flow_id in data with key 'flow_id'")
	return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
}

func (self *ClientIG[C, D]) Update(typeUpdate string, data ...map[string]any) gralresponse.Responser {
	dataParam := make(map[string]any)

	if len(data) > 0 {
		dataParam = data[0]
	}

	switch typeUpdate {
	case ig.IG_UPDATE_WELCOME_MESSAGE_FLOWS:
		return self.updateWelcomeMessageFlowsADS(dataParam)
	default:
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("typeUpdate not recognized. typeUpdate expect 'welcome_message_flows'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}
