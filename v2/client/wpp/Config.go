package wpp

import (
	"net/http"
	"path"

	response "github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type Config struct {
	Token      string `json:"token"`
	wabaID     string
	phoneID    string
	businessID string
	// Path del archivo .env incluyendo el nombre del archivo sin la extensión ej: file: /.../../config_env.env -> EnvFilePath: /.../../config_env
	EnvFilePath string `json:"env_file_path"`
	Error       error
	// Url del servidor WebHook con ruta /ws para conectar con el servidor WebSocket ej: wss://path.com/ws if not set send data to client.Broadcast(data)
	WebhookSocket string    `json:"webhook_socket"`
	EventHandle   func(any) // Funcion para manejar los eventos del servidor WebHook WebSocket
	pathVersion   string    // URL, ej: https://graph.facebook.com/v15.0
	pathPhone     string    // URL, ej: https://graph.facebook.com/v15.0/PHONE_NUMBER_ID
	pathBusiness  string    // URL, ej: https://graph.facebook.com/v15.0/BUSINESS_ACCOUNT_ID
	pathWaba      string    // URL, ej: https://graph.facebook.com/v15.0/WABA_ID
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

type Options func(*Config)

func WithToken(token string) Options {
	return func(c *Config) {
		c.Token = token
	}
}

func WithWabaID(waba_id string) Options {
	return func(c *Config) {
		c.wabaID = waba_id
		// c.pathWaba = path.Join(c.cLOUD_API_VERSION, c.wabaID)
	}
}

func WithBusinessID(business_id string) Options {
	return func(c *Config) {
		c.businessID = business_id
		// c.pathBusiness = path.Join(c.cLOUD_API_VERSION, c.businessID)
	}
}

func WithPhoneID(phone_id string) Options {
	return func(c *Config) {
		c.phoneID = phone_id
		// c.pathPhone = path.Join(c.cLOUD_API_VERSION, c.phoneID)
	}
}

func WithEnvFilePath(env_file_path string) Options {
	return func(c *Config) {
		c.EnvFilePath = env_file_path
	}
}

func WithWebhookSocket(webhook_socket string) Options {
	return func(c *Config) {
		c.WebhookSocket = webhook_socket
	}
}

func WithEventHandle(event_handle func(any)) Options {
	return func(c *Config) {
		c.EventHandle = event_handle
	}
}

// setter
func (c *Config) setBaseUrl(wa_base_url string) {
	c.wA_BASE_URL = wa_base_url
}

func (c *Config) setM4DAppId(m4d_app_id string) {
	c.m4D_APP_ID = m4d_app_id
}

func (c *Config) setM4DAppSecret(m4d_app_secret string) {
	c.m4D_APP_SECRET = m4d_app_secret
}

func (c *Config) setPhoneID(phone_id string) {
	c.phoneID = phone_id
	c.pathPhone = path.Join(c.getVersion(), c.phoneID)
}

func (c *Config) setWabaID(waba_id string) {
	c.wabaID = waba_id
	c.pathWaba = path.Join(c.getVersion(), c.wabaID)
}

func (c *Config) setBusinessID(business_id string) {
	c.businessID = business_id
	c.pathBusiness = path.Join(c.getVersion(), c.businessID)
}

func (c *Config) setCloudApiAccessToken(cloud_api_access_token string) {
	c.cLOUD_API_ACCESS_TOKEN = cloud_api_access_token
}

func (c *Config) setCloudApiVersion(cloud_api_version string) {
	c.cLOUD_API_VERSION = cloud_api_version
	// c.pathVersion = path.Join(c.cLOUD_API_VERSION)
	// c.pathPhone = path.Join(c.cLOUD_API_VERSION, c.phoneID)
	// c.pathBusiness = path.Join(c.cLOUD_API_VERSION, c.businessID)
	// c.pathWaba = path.Join(c.cLOUD_API_VERSION, c.wabaID)
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
func (c Config) getWabaID() string {
	return c.wabaID
}

func (c Config) getBusinessID() string {
	return c.businessID
}

func (c Config) getphoneID() string {
	return c.phoneID
}

func (c Config) getBaseUrl() string {
	return c.wA_BASE_URL
}

func (c Config) getVersion() string {
	return c.cLOUD_API_VERSION
}
