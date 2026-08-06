package embeddedsignup

// TypesEvents embedded signup
type EVENT_TYPE string

const (
	//  indica que se completó correctamente el proceso de la API de la nube.
	FINISH EVENT_TYPE = "FINISH"
	//  indica que el usuario completó el proceso sin un número de teléfono.
	FINISH_ONLY_WABA EVENT_TYPE = "FINISH_ONLY_WABA"
	// indica que el usuario completó el proceso con un número de la app de WhatsApp Business.
	FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING EVENT_TYPE = "FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING"
	CANCEL                                  EVENT_TYPE = "CANCEL"
	ERROR                                   EVENT_TYPE = "ERROR"
	FINISH_OBO_MIGRATION                    EVENT_TYPE = "FINISH_OBO_MIGRATION"
	FINISH_GRANT_ONLY_API_ACCESS            EVENT_TYPE = "FINISH_GRANT_ONLY_API_ACCESS"
)

func (e EVENT_TYPE) String() string {
	if val, ok := map[EVENT_TYPE]string{
		CANCEL:                                  "CANCEL",
		ERROR:                                   "ERROR",
		FINISH:                                  "FINISH",
		FINISH_ONLY_WABA:                        "FINISH_ONLY_WABA",
		FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING: "FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING",
		FINISH_OBO_MIGRATION:                    "FINISH_OBO_MIGRATION",
		FINISH_GRANT_ONLY_API_ACCESS:            "FINISH_GRANT_ONLY_API_ACCESS",
	}[e]; ok {
		return val
	}

	return string(e)
}

type RegisterSessionEvent struct {
	Data WAEmbeddedSignupData `json:"data" form:"data"`
	Type string               `json:"type" form:"type"`
	/*
		FINISH: Indicate that the Cloud Api process completed successfully.
		FINISH_ONLY_WABA: Indicate that the user completed successfully the process without a phone number.
		FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING: Indicate the the user completed successfully the process with a phone number of WhatsApp Business app.
		FINISH_OBO_MIGRATION: Indicate that the user completed the representation process
		FINISH_GRANT_ONLY_API_ACCESS: Indicate that the user completed that grants api-only access
		ERROR: Indicate that the user encountered an error during the process
		CANCEL: Indicate that the user canceled the process
	*/
	Event EVENT_TYPE `json:"event" form:"event"`
}

type WAEmbeddedSignupData struct {
	PhoneNumberID       string   `json:"phone_number_id,omitempty" form:"phone_number_id"`
	WabaID              string   `json:"waba_id,omitempty" form:"waba_id"`
	BusinessID          string   `json:"business_id,omitempty" form:"business_id"`
	AdAccountIDs        []string `json:"ad_account_ids,omitempty" form:"ad_account_ids"`               // only included if customer selected ad accounts
	PageIDs             []string `json:"page_ids,omitempty" form:"page_ids"`                           // only included if customer selected Facebook Pages
	DatasetIDs          []string `json:"dataset_ids,omitempty" form:"dataset_ids"`                     // only included if customer selected datasets
	CatalogIds          []string `json:"catalog_ids,omitempty" form:"catalog_ids"`                     // only included if customer selected catalogs
	InstagramAccountIds []string `json:"instagram_account_ids,omitempty" form:"instagram_account_ids"` // only included if customer selected Instagram accounts
	WabaIds             []string `json:"waba_ids,omitempty" form:"waba_ids"`                           // only included for multi-WABA flows
	CurrentStep         string   `json:"current_step,omitempty" form:"current_step"`                   // indicate which screen the business customer was on when they abandoned the process, example: current_step: 'PHONE_NUMBER_SETUP'
	ErrorMessage        string   `json:"error_message,omitempty" form:"error_message"`
	ErrorCode           string   `json:"error_code,omitempty" form:"error_code"`
	ErrorId             string   `json:"error_id,omitempty" form:"error_id"`
	SessionId           string   `json:"session_id,omitempty" form:"session_id"`
	Timestamp           string   `json:"timestamp,omitempty" form:"timestamp"`
}
