package clientoficial

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

var (
	WA_BASE_URL                string
	M4D_APP_ID                 string
	M4D_APP_SECRET             string
	WA_PHONE_NUMBER_ID         string
	WA_BUSINESS_ACCOUNT_ID     string
	CLOUD_API_ACCESS_TOKEN     string
	CLOUD_API_VERSION          string
	WEBHOOK_ENDPOINT           string
	WEBHOOK_VERIFICATION_TOKEN string
	LISTENER_PORT              string
	DEBUG                      string
	MAX_RETRIES_AFTER_WAIT     string
	REQUEST_TIMEOUT            string
)

func setEnv(envPath string) {
	pathDir := path.Dir(envPath)
	envName := path.Base(envPath)
	viper.AddConfigPath(pathDir)
	viper.SetConfigType("env")
	// viper.SetConfigName("config_env.env")
	viper.SetConfigName(fmt.Sprintf("%s.env", envName))
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("\033[31mError: No encontrado archivo app.env ni .cobraToml de tipo (toml) en\033[30m %s\n", pathDir)
		os.Exit(2)
	} else {
		// Variables .env API_WHATSAPP
		WA_BASE_URL = viper.GetString("WA_BASE_URL")
		M4D_APP_ID = viper.GetString("M4D_APP_ID")
		M4D_APP_SECRET = viper.GetString("M4D_APP_SECRET")
		WA_PHONE_NUMBER_ID = viper.GetString("WA_PHONE_NUMBER_ID")
		WA_BUSINESS_ACCOUNT_ID = viper.GetString("WA_BUSINESS_ACCOUNT_ID")
		CLOUD_API_ACCESS_TOKEN = viper.GetString("CLOUD_API_ACCESS_TOKEN")
		CLOUD_API_VERSION = viper.GetString("CLOUD_API_VERSION")
		WEBHOOK_ENDPOINT = viper.GetString("WEBHOOK_ENDPOINT")
		WEBHOOK_VERIFICATION_TOKEN = viper.GetString("WEBHOOK_VERIFICATION_TOKEN")
		LISTENER_PORT = viper.GetString("LISTENER_PORT")
		DEBUG = viper.GetString("DEBUG")
		MAX_RETRIES_AFTER_WAIT = viper.GetString("MAX_RETRIES_AFTER_WAIT")
		REQUEST_TIMEOUT = viper.GetString("REQUEST_TIMEOUT")
	}
}
