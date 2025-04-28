//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"fmt"
	"path"

	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
	"github.com/spf13/viper"
)

var (
// WA_BASE_URL                string
// M4D_APP_ID                 string
// M4D_APP_SECRET             string
// WA_PHONE_NUMBER_ID         string
// WA_BUSINESS_ACCOUNT_ID     string
// CLOUD_API_ACCESS_TOKEN     string
// CLOUD_API_VERSION          string
// WEBHOOK_ENDPOINT           string
// WEBHOOK_VERIFICATION_TOKEN string
// LISTENER_PORT              string
// DEBUG                      string
// MAX_RETRIES_AFTER_WAIT     string
// REQUEST_TIMEOUT            string
// TOKEN                      string
)

func setEnv(c *Config) error {
	var envPath string = c.EnvFilePath
	pathDir := path.Dir(envPath)
	envName := path.Base(envPath)
	viper.AddConfigPath(pathDir)
	viper.SetConfigType("env")
	// viper.SetConfigName("config_env.env")
	viper.SetConfigName(fmt.Sprintf("%s.env", envName))
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("\033[31mError: No encontrado archivo app.env ni .cobraToml de tipo (toml) en\033[30m %s\n", pathDir)
		return response.NewError(&response.Error{
			Type:    types.TypeErrorConfig,
			Code:    types.CodeErrorEnvNotFound,
			Message: types.MsgErrorEnvNotFound,
		})
	} else {
		// Variables .env API_WHATSAPP
		c.wA_BASE_URL = viper.GetString("WA_BASE_URL")
		c.m4D_APP_ID = viper.GetString("M4D_APP_ID")
		c.m4D_APP_SECRET = viper.GetString("M4D_APP_SECRET")
		if c.WaBusinessAccountId == "" {
			c.wA_BUSINESS_ACCOUNT_ID = viper.GetString("WA_BUSINESS_ACCOUNT_ID")
		} else {
			c.wA_BUSINESS_ACCOUNT_ID = c.WaBusinessAccountId
		}

		if c.WaPhoneNumberId == "" {
			c.wA_PHONE_NUMBER_ID = viper.GetString("WA_PHONE_NUMBER_ID")
		} else {
			c.wA_PHONE_NUMBER_ID = c.WaPhoneNumberId
		}

		if c.Token == "" {
			c.Token = viper.GetString("CLOUD_API_ACCESS_TOKEN")
			c.cLOUD_API_ACCESS_TOKEN = c.Token
		}
		c.cLOUD_API_VERSION = viper.GetString("CLOUD_API_VERSION")
		c.wEBHOOK_ENDPOINT = viper.GetString("WEBHOOK_ENDPOINT")
		c.wEBHOOK_VERIFICATION_TOKEN = viper.GetString("WEBHOOK_VERIFICATION_TOKEN")
		c.lISTENER_PORT = viper.GetString("LISTENER_PORT")
		c.dEBUG = viper.GetString("DEBUG")
		c.mAX_RETRIES_AFTER_WAIT = viper.GetString("MAX_RETRIES_AFTER_WAIT")
		c.rEQUEST_TIMEOUT = viper.GetString("REQUEST_TIMEOUT")
	}
	return nil
}
