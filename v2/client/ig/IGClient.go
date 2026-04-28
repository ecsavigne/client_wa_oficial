package ig

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ecsavigne/client_wa_oficial/v2/client"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	generalpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	gralpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	gralrequest "github.com/ecsavigne/client_wa_oficial/v2/types/general/request"
	gralresponse "github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	"github.com/ecsavigne/client_wa_oficial/v2/types/ig"
	igpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/ig/gen/igpb/v1"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type Config struct {
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
}

func (self *Config) GetType() client.TYPE_CONFIG          { return client.TYPE_CONFIG_IG }
func (self *Config) SetToken(token string)                { self.token = token }
func (self *Config) GetToken() string                     { return self.token }
func (self *Config) SetUserID(userID string)              { self.userID = userID }
func (self Config) GetUserID() string                     { return self.userID }
func (self *Config) SetEnvFilePath(envFilePath string)    { self.envFilePath = envFilePath }
func (self Config) GetEnvFilePath() string                { return self.envFilePath }
func (self *Config) SetEventHandle(eventHandle func(any)) { self.eventHandle = eventHandle }
func (self Config) GetEventHandle() func(any)             { return self.eventHandle }
func (self *Config) SetVersion(version string)            { self.version = version }
func (self *Config) GetVersion() string                   { return self.version }
func (self *Config) SetBaseUrlFacebook(baseUrlFacebook string) {
	self.baseUrlFacebook = baseUrlFacebook
}
func (self Config) GetBaseUrlFacebook() string     { return self.baseUrlFacebook }
func (self *Config) SetBaseUrlIG(baseUrlIG string) { self.baseUrlIG = baseUrlIG }
func (self Config) GetBaseUrlIG() string           { return self.baseUrlIG }
func (self *Config) SetBaseUrl(baseUrl string)     { self.baseUrlIG = baseUrl }
func (self *Config) GetBaseUrl() string            { return self.baseUrlIG }
func (self *Config) SetAppId(appID string)         { self.appId = appID }
func (self Config) GetAppId() string               { return self.appId }
func (self *Config) SetAppSecret(appSecret string) { self.appSecret = appSecret }
func (self Config) GetAppSecret() string           { return self.appSecret }
func (self Config) String() string {
	return fmt.Sprintf("Config{token: %s, userID: %s, envFilePath: %s, version: %s, baseUrlFacebook: %s, baseUrlIG: %s, appId: %s, appSecret: %s}",
		self.token, self.userID, self.envFilePath, self.version, self.baseUrlFacebook, self.baseUrlIG, self.appId, self.appSecret)
}

type OptionsClientIG func(*Config) // Opts

func WithToken(token string) OptionsClientIG {
	return func(c *Config) {
		c.token = token
	}
}

func WithUserID(userID string) OptionsClientIG {
	return func(c *Config) {
		c.userID = userID
	}
}

func WithEnvFilePath(envFilePath string) OptionsClientIG {
	return func(c *Config) {
		c.envFilePath = envFilePath
	}
}

func WithEventHandle(eventHandle func(any)) OptionsClientIG {
	return func(c *Config) {
		c.eventHandle = eventHandle
	}
}

func WithVersion(version string) OptionsClientIG {
	return func(c *Config) {
		c.version = version
	}
}

func WithBaseUrlFacebook(baseUrl string) OptionsClientIG {
	return func(c *Config) {
		c.baseUrlFacebook = baseUrl
	}
}

func WithBaseUrlIG(baseUrl string) OptionsClientIG {
	return func(c *Config) {
		c.baseUrlIG = baseUrl
	}
}

func WithAppId(appID string) OptionsClientIG {
	return func(c *Config) {
		c.appId = appID
	}
}

func WithAppSecret(appSecret string) OptionsClientIG {
	return func(c *Config) {
		c.appSecret = appSecret
	}
}

func setEnv(c *Config) error {
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

func defaultConfig() *Config {
	c := &Config{
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

type ClientIG struct {
	config     *Config
	typeClient client.CLIENT_TYPE
	url        *url.URL
}

func (self *ClientIG) String() string {
	return fmt.Sprintf("ClientIG{config: %s, typeClient: %s, url_base: %s}", self.config.String(), self.typeClient.String(), self.GetConfig().GetBaseUrl())
}

func (self *ClientIG) MultipartRequest(method string, data proto.Message, ePoint string) (*http.Request, error) {
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

func NewClientIG(opts ...OptionsClientIG) (*ClientIG, error) {
	c := defaultConfig()

	for _, opt := range opts {
		opt(c)
	}

	err := setEnv(c)
	if err != nil {
		return nil, err
	}

	cl := &ClientIG{
		config:     c,
		typeClient: client.CLIENT_IG,
	}

	return cl, nil
}

// func (self *ClientIG) createUrl() {
// 	// return self.typeClient.String()
// }

func (self *ClientIG) GetType() string                { return self.typeClient.String() }
func (self *ClientIG) GetConfig() client.ConfigClient { return self.config }

// func getPossivelResponse() gralresponse.ResponseType{}

func (self *ClientIG) executeRequest(method string, ePoint string, data any, isID bool, resp_ ...gralresponse.ResponseType) gralresponse.Responser {
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

func (self *ClientIG) sendTextMessage(msg *igpbv1.InstagramTextMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
}

func (self *ClientIG) sendReactionMessage(msg *igpbv1.InstagramReactionMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
}

// upload media to Meta servers and get media_id if msg.GetFileHeader() exist, upload media to Meta servers and get media_id, then send media message with media_id
// if msg.GetFileHeader() == nil and exist msg.GetMessageLink() or msg.GetMessage().GetId() != "" then send media message with message link or media id
// if msg.GetFileHeader() == nil and msg.GetFileHeader() == nil and not exist msg.GetMessageLink() and msg.GetMessage().GetId() == "" then return error response with message "Media file is required for media message"() == "" and msg.GetMessage().GetId() == "" then return error response with message "Media file not found. Please provide a media file or a media link or a media id"
func (self *ClientIG) sendMediaMessage(msg *igpbv1.InstagramMediaMessage) gralresponse.Responser {
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

	if msg.GetMessage().GetAttachment().GetType() == ig.IG_MEDIA_MESSAGE_TYPE_LIKE_HEART || msg.GetMessage().GetAttachments() != nil {
		return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
	}

	// if msg.GetFileHeader() != nil {

	return nil
}

func (self *ClientIG) sendMediaShareMessage(msg *igpbv1.InstagramMediaShareMessage) gralresponse.Responser {
	return nil
}

/**
* @name: sendInstagramQuickRepliesMessage
* @description: Send Instagram Quick Replies message to Instagram
* @param {proto.Message} msg e.g: *igpbv1.InstagramQuickRepliesMessage
* @return gralresponse.Responser
 */
func (self *ClientIG) sendInstagramQuickRepliesMessage(msg *igpbv1.InstagramQuickRepliesMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messages", msg, true, gralresponse.SentMessageResponse)
}

func (self *ClientIG) sendInstagramPersistentMenu(msg *igpbv1.InstagramPersistentMenuMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messenger_profile", msg, true)
}

/**
* @name: SendAction
* @description: Send action to instagram
* @param {string} scope_id the id of the user to send the action to
* @param {string} action the action to send. Can be "typing_on" or "typing_off"
* @return gralresponse.Responser
 */
func (self *ClientIG) sendAction(scope_id, action string) gralresponse.Responser {
	if action != "typing_on" && action != "typing_off" {
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Action not recognized. Action expect 'typing_on' or 'typing_off'")
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

func (self *ClientIG) SendPresence(recipient_id, action string) gralresponse.Responser {
	return self.sendAction(recipient_id, action)
}

func (self *ClientIG) sendInstagramIceBreakersMessage(msg *igpbv1.InstagramIceBreakersMessage) gralresponse.Responser {
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
func (self *ClientIG) createInstagramWelcomeMessageFlows(msg *igpbv1.InstagramWelcomeMessageFlows) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, fmt.Sprintf("/%s/%s", "me", "welcome_message_flows"), msg, false)
}

/**
* @name: SendMessage
* @description: Send message to instagram
* @param {proto.Message} msg e.g: *igpbv1.InstagramTextMessage, *igpbv1.InstagramMediaMessage, *igpbv1.InstagramMediaShareMessage
* @return gralresponse.Responser
 */
func (self *ClientIG) SendMessage(msg proto.Message) gralresponse.Responser {
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
	default:
		fmt.Printf("Error in SendMessage, file: IGClient.go. Message type not recognized. Type of message is: %T\n", v)
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Message not recognized. Message.type expect 'text', 'audio', 'image', 'video', 'document', 'sticker', 'location', 'contact', 'template', 'interactive', 'reaction'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}

// Gets
func (self *ClientIG) getInfoAccountBusiness() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/", types.QueryData{
		"fields": "id,user_id,media_count,name,username,followers_count,follows_count,profile_picture_url",
	}, true, gralresponse.InfoAccountBusinessResponse)
}

func (self *ClientIG) getInstagramPersistentMenu() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/messenger_profile", types.QueryData{
		"fields": "persistent_menu",
	}, true)
}

func (self *ClientIG) getInstagramIceBreakers() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/messenger_profile", types.QueryData{
		"fields": "ice_breakers",
	}, true)
}

func (self *ClientIG) getInstagramLink() gralresponse.Responser {
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
	unknow.SetData(map[string]*structpb.Value{
		"link": structpb.NewStringValue(link),
	})

	return gralresponse.NewResponse(unknow, gralresponse.ResponseUnknow)
}

func (self *ClientIG) getWelcomeMessageFlowsADS() gralresponse.Responser {
	return self.executeRequest(http.MethodGet, "/me/welcome_message_flows", types.QueryData{}, false)
}

func (self *ClientIG) Get(type_info string) gralresponse.Responser {
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
	default:
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("type_info not recognized. type_info expect 'account_business', 'persistent_menu', 'ice_breakers', 'link', 'welcome_message_flows'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}

// Deletes
func (self *ClientIG) deletePersistentMenu() gralresponse.Responser {
	return self.executeRequest(http.MethodDelete, "/messenger_profile", types.QueryData{
		"fields": []string{"persistent_menu"},
	}, true)
}

func (self *ClientIG) deleteIceBreakers() gralresponse.Responser {
	return self.executeRequest(http.MethodDelete, "/messenger_profile", types.QueryData{
		"fields": []string{"ice_breakers"},
	}, true)
}

func (self *ClientIG) deleteWelcomeMessageFlowsADS(data types.QueryData) gralresponse.Responser {
	if id, ok := data["flow_id"]; ok && id != "" {
		id = fmt.Sprintf("%c%s", '?', data.String())
		return self.executeRequest(http.MethodDelete, fmt.Sprintf("/me/welcome_message_flows%s", id), types.QueryData{}, false)
	}

	errorResponse := &gralpbv1.ResponseError{}
	errorResponse.SetCode(401)
	errorResponse.SetMessage("flow_id is required to delete welcome message flow. Please provide flow_id in data with key 'flow_id'")
	return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
}

func (self *ClientIG) Delete(typeDelete string, data ...map[string]any) gralresponse.Responser {
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
	default:
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("typeDelete not recognized. typeDelete expect 'persistent_menu', 'ice_breakers', 'welcome_message_flows'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}

// Updates
func (self *ClientIG) updateWelcomeMessageFlowsADS(data types.QueryData) gralresponse.Responser {
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

func (self *ClientIG) Update(typeUpdate string, data ...map[string]any) gralresponse.Responser {
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
