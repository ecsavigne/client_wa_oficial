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
		c.setBaseUrl(viper.GetString("WA_BASE_URL"))
		c.setM4DAppId(viper.GetString("M4D_APP_ID"))
		c.setM4DAppSecret(viper.GetString("M4D_APP_SECRET"))
		if c.getWabaID() == "" {
			// c.WaBusinessAccountId = viper.GetString("WA_BUSINESS_ACCOUNT_ID")
			c.setWabaID(viper.GetString("WA_BUSINESS_ACCOUNT_ID"))
		}

		if c.getphoneID() == "" {
			c.setPhoneID(viper.GetString("WA_PHONE_NUMBER_ID"))
		}

		if c.Token == "" {
			c.Token = viper.GetString("CLOUD_API_ACCESS_TOKEN")
			c.cLOUD_API_ACCESS_TOKEN = c.Token
		}
		c.setCloudApiVersion(viper.GetString("CLOUD_API_VERSION"))
		c.setWebhookEndpoint(viper.GetString("WEBHOOK_ENDPOINT"))
		c.setWebhookVerificationToken(viper.GetString("WEBHOOK_VERIFICATION_TOKEN"))
		c.setListenerPort(viper.GetString("LISTENER_PORT"))
		c.setDebug(viper.GetString("DEBUG"))
		c.setMaxRetriesAfterWait(viper.GetString("MAX_RETRIES_AFTER_WAIT"))
		c.setRequestTimeout(viper.GetString("REQUEST_TIMEOUT"))
	}
	return nil
}
