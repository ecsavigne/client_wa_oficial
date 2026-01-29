package embeddedsignup

// TypesEvents embedded signup
type EVENT_TYPE string
type FLOW_FINISH_TYPE EVENT_TYPE

const (
	FINISH                                  FLOW_FINISH_TYPE = "FINISH"
	FINISH_ONLY_WABA                        FLOW_FINISH_TYPE = "FINISH_ONLY_WABA"
	FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING FLOW_FINISH_TYPE = "FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING"
	CANCEL                                  EVENT_TYPE       = "CANCEL"
)

func (e EVENT_TYPE) String() string {
	if val, ok := map[EVENT_TYPE]string{
		CANCEL: "CANCEL",
	}[e]; ok {
		return val
	} else if val, ok := map[FLOW_FINISH_TYPE]string{
		FINISH:                                  "FINISH",
		FINISH_ONLY_WABA:                        "FINISH_ONLY_WABA",
		FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING: "FINISH_WHATSAPP_BUSINESS_APP_ONBOARDING",
	}[FLOW_FINISH_TYPE(e)]; ok {
		return val
	}

	return string(e)
}

type RegisterSessionEvent struct {
	Data  WAEmbeddedSignupData `json:"data" form:"data"`
	Type  string               `json:"type" form:"type"`
	Event EVENT_TYPE           `json:"event" form:"event"`
}

type WAEmbeddedSignupData struct {
	PhoneNumberID string   `json:"phone_number_id,omitempty" form:"phone_number_id"`
	WabaID        string   `json:"waba_id,omitempty" form:"waba_id"`
	BusinessID    string   `json:"business_id,omitempty" form:"business_id"`
	AdAccountIDs  []string `json:"ad_account_ids,omitempty" form:"ad_account_ids"`
	PageIDs       []string `json:"page_ids,omitempty" form:"page_ids"`
	DatasetIDs    []string `json:"dataset_ids,omitempty" form:"dataset_ids"`
	CurrentStep   string   `json:"current_step,omitempty" form:"current_step"`
	ErrorMessage  string   `json:"error_message,omitempty" form:"error_message"`
	ErrorId       string   `json:"error_id,omitempty" form:"error_id"`
	SessionId     string   `json:"session_id,omitempty" form:"session_id"`
	Timestamp     string   `json:"timestamp,omitempty" form:"timestamp"`
}
