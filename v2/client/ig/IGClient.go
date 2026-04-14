package ig

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ecsavigne/client_wa_oficial/v2/client"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	gralpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	gralrequest "github.com/ecsavigne/client_wa_oficial/v2/types/general/request"
	gralresponse "github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	igpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/ig/gen/igpb/v1"
	igrsponse "github.com/ecsavigne/client_wa_oficial/v2/types/ig/response"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
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

func (self *ClientIG) createUrl() {
	// return self.typeClient.String()
}

func (self *ClientIG) GetType() string                { return self.typeClient.String() }
func (self *ClientIG) GetConfig() client.ConfigClient { return self.config }

// func getPossivelResponse() gralresponse.ResponseType{}

func (self *ClientIG) executeRequest(methoth string, ePoint string, data any, resp_ ...gralresponse.ResponseType) gralresponse.Responser {
	respType := gralresponse.ResponseUnknow
	if len(resp_) > 0 {
		respType = resp_[0]
	}

	// if msg.GetMessageLink() != "" || msg.GetFileHeader() != nil {
	// 	multipartRequest(methoth, ePoint, c.Config, msg)
	// } else {
	req, err := gralrequest.DefaultRequest(self.GetConfig(), methoth, ePoint, data)
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

func (self *ClientIG) SendTextMessage(msg *igpbv1.InstagramTextMessage) gralresponse.Responser {
	return self.executeRequest(http.MethodPost, "/messages", msg, igrsponse.ResponseSentMessage)
}

// upload media to Meta servers and get media_id if msg.GetFileHeader() exist, upload media to Meta servers and get media_id, then send media message with media_id
// if msg.GetFileHeader() == nil and exist msg.GetMessageLink() or msg.GetMessage().GetId() != "" then send media message with message link or media id
// if msg.GetFileHeader() == nil and msg.GetFileHeader() == nil and not exist msg.GetMessageLink() and msg.GetMessage().GetId() == "" then return error response with message "Media file is required for media message"() == "" and msg.GetMessage().GetId() == "" then return error response with message "Media file not found. Please provide a media file or a media link or a media id"
func (self *ClientIG) SendMediaMessage(msg *igpbv1.InstagramMediaMessage) gralresponse.Responser {
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
		return self.executeRequest(http.MethodPost, "/messages", msg, igrsponse.ResponseSentMessage)
	}

	// if send media id media message
	if msg.GetMessage().GetAttachment().GetPayload().GetAttachmentId() != "" {
		msg.SetFileDescriptor(nil)
		msg.GetMessage().GetAttachment().GetPayload().SetUrl("")
		return self.executeRequest(http.MethodPost, "/messages", msg, igrsponse.ResponseSentMessage)
	}

	// if msg.GetFileHeader() != nil {

	return nil
}

func (self *ClientIG) SendMediaShareMessage(msg *igpbv1.InstagramMediaShareMessage) gralresponse.Responser {
	return nil
}

// func (c *ClientIG) SendMessage(m message.Messager) response.Responser {
// 	switch m.GetType() {
// 	case wpp.MessageTypeAudio:
// 		return c.sendAudioMessage(m)
// 	case wpp.MessageTypeContact:
// 		return c.sendContactMessage(m)
// 	case wpp.MessageTypeDocument:
// 		return c.sendDocumentMessage(m)
// 	case wpp.MessageTypeImage:
// 		return c.sendImageMessage(m)
// 	case wpp.MessageTypeInteractive:
// 		return c.sendInteractive(m)
// 	case wpp.MessageTypeLocation:
// 		return c.sendLocationMessage(m)
// 	case wpp.MessageTypeReaction:
// 		return c.sendReaction(m)
// 	case wpp.MessageTypeResponse:
// 		return c.sendResponseMsg(m)
// 	case wpp.MessageTypeSticker:
// 		return c.sendStickerMessage(m)
// 	case wpp.MessageTypeTemplate:
// 		return c.sendTemplate(m)
// 	case wpp.MessageTypeText:
// 		return c.sendTextMessage(m)
// 	case wpp.MessageTypeVideo:
// 		return c.sendVideoMessage(m)
// 	default:
// 		return response.NewError(&response.Error{
// 			Type:    response.ResponseError,
// 			Code:    401,
// 			Message: fmt.Sprintf("Message not recognized. Message.type expect '%v'", []string{"text", "audio", "image", "video", "document", "sticker", "location", "contact", "template", "interactive", "reaction"}),
// 		})
// 	}
// }

func (self *ClientIG) SendMessage(msg proto.Message) gralresponse.Responser {
	switch v := msg.(type) {
	case *igpbv1.InstagramTextMessage:
		return self.SendTextMessage(v)
	case *igpbv1.InstagramMediaMessage:
		return self.SendMediaMessage(v)
	case *igpbv1.InstagramMediaShareMessage:
		return self.SendMediaShareMessage(v)
	default:
		errorResponse := &gralpbv1.ResponseError{}
		errorResponse.SetCode(401)
		errorResponse.SetMessage("Message not recognized. Message.type expect 'text', 'audio', 'image', 'video', 'document', 'sticker', 'location', 'contact', 'template', 'interactive', 'reaction'")
		return gralresponse.NewResponse(errorResponse, gralresponse.ResponseError)
	}
}
