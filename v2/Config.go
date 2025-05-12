package clientoficial

import (
	"net/http"
	"path"

	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

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
	path          string    // URL, ej: https://graph.facebook.com/v15.0/PHONE_NUMBER_ID
	pathVersion   string    // URL, ej: https://graph.facebook.com/v15.0
	pathBusiness  string    // URL, ej: https://graph.facebook.com/v15.0/BUSINESS_ACCOUNT_ID
	clientHttp
	request        *http.Request
	MediaInfo      *response.MediaInfo
	wA_BASE_URL    string
	m4D_APP_ID     string
	m4D_APP_SECRET string
	// wA_PHONE_NUMBER_ID         string
	// wA_BUSINESS_ACCOUNT_ID     string
	cLOUD_API_ACCESS_TOKEN     string
	cLOUD_API_VERSION          string
	wEBHOOK_ENDPOINT           string
	wEBHOOK_VERIFICATION_TOKEN string
	lISTENER_PORT              string
	dEBUG                      string
	mAX_RETRIES_AFTER_WAIT     string
	rEQUEST_TIMEOUT            string
	// tOKEN                      string
}

func (c *Config) setWaBaseUrl(wa_base_url string) {
	c.wA_BASE_URL = wa_base_url
}

func (c *Config) setM4DAppId(m4d_app_id string) {
	c.m4D_APP_ID = m4d_app_id
}

func (c *Config) setM4DAppSecret(m4d_app_secret string) {
	c.m4D_APP_SECRET = m4d_app_secret
}

func (c *Config) setWaPhoneNumberId(wa_phone_number_id string) {
	c.WaPhoneNumberId = wa_phone_number_id
	c.path = path.Join(c.cLOUD_API_VERSION, c.WaPhoneNumberId)
}

func (c *Config) setWaBusinessAccountId(wa_business_account_id string) {
	c.WaBusinessAccountId = wa_business_account_id
	c.pathBusiness = path.Join(c.cLOUD_API_VERSION, c.WaBusinessAccountId)
}

func (c *Config) setCloudApiAccessToken(cloud_api_access_token string) {
	c.cLOUD_API_ACCESS_TOKEN = cloud_api_access_token
}

func (c *Config) setCloudApiVersion(cloud_api_version string) {
	c.cLOUD_API_VERSION = cloud_api_version
	c.pathVersion = path.Join(c.cLOUD_API_VERSION)
	c.path = path.Join(c.cLOUD_API_VERSION, c.WaPhoneNumberId)
	c.pathBusiness = path.Join(c.cLOUD_API_VERSION, c.WaBusinessAccountId)
}

func (c *Config) setWebhookEndpoint(webhook_endpoint string) {
	c.wEBHOOK_ENDPOINT = webhook_endpoint
}

func (c *Config) setWebhookVerificationToken(webhook_verification_token string) {
	c.wEBHOOK_VERIFICATION_TOKEN = webhook_verification_token
}

func (c *Config) setListenerPort(listener_port string) {
	c.lISTENER_PORT = listener_port
}

func (c *Config) setDebug(debug string) {
	c.dEBUG = debug
}

func (c *Config) setMaxRetriesAfterWait(max_retries_after_wait string) {
	c.mAX_RETRIES_AFTER_WAIT = max_retries_after_wait
}

func (c *Config) setRequestTimeout(request_timeout string) {
	c.rEQUEST_TIMEOUT = request_timeout
}

// getter
