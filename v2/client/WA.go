//lint:file-ignore ST1005 Ignore capitalized strings error
package clientoficial

import (
	"fmt"
	"path"

	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
	"github.com/spf13/viper"
)

func setEnv(c *Config) error {
	var envPath string = c.EnvFilePath
	pathDir := path.Dir(envPath)
	envName := path.Base(envPath)
	viper.AddConfigPath(pathDir)
	viper.SetConfigType("env")
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
			c.WaBusinessAccountId = viper.GetString("WA_BUSINESS_ACCOUNT_ID")
		}

		if c.WaPhoneNumberId == "" {
			c.WaPhoneNumberId = viper.GetString("WA_PHONE_NUMBER_ID")
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
