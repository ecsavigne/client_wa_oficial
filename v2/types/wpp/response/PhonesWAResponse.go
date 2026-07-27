package response

/*
Represent the response from the GET "https://graph.facebook.com/v22.0/{whatsapp-business-account-id}/phone_numbers
*/

type PhonesWA struct {
	KernelResponser
	Type         string      `json:"type,omitempty"`
	ResponseType string      `json:"response_type,omitempty"`
	Data         []PhoneInfo `json:"data"`
	Paging       *Paging     `json:"paging,omitempty"`
}

type PhoneInfo struct {
	VerifiedName              string                `json:"verified_name,omitempty"`
	CodeVerificationStatus    string                `json:"code_verification_status,omitempty"`
	DisplayPhoneNumber        string                `json:"display_phone_number,omitempty"`
	QualityRating             string                `json:"quality_rating,omitempty"`
	PlatformType              string                `json:"platform_type,omitempty"`
	ID                        string                `json:"id,omitempty"`
	NameStatus                string                `json:"name_status,omitempty"`
	AccountMode               string                `json:"account_mode,omitempty"`
	IsOfficialBusinessAccount bool                  `json:"is_official_business_account,omitempty"`
	LastOnboardedTime         string                `json:"last_onboarded_time"`
	Status                    string                `json:"status,omitempty"`
	IsOnBizApp                bool                  `json:"is_on_biz_app,omitempty"`
	IsPreverifiedNumber       bool                  `json:"is_preverified_number,omitempty"`
	NewNameStatus             string                `json:"new_name_status,omitempty"`
	Throughput                *Throughput           `json:"throughput,omitempty"`
	WebhookConfiguration      *WebhookConfiguration `json:"webhook_configuration,omitempty"`
}

type Throughput struct {
	Level string `json:"level,omitempty"`
}

type WebhookConfiguration struct {
	Application string `json:"application,omitempty"`
}

type Paging struct {
	Cursors  *Cursors `json:"cursors,omitempty"`
	Next     string   `json:"next,omitempty"`
	Previous string   `json:"previous,omitempty"`
}

type Cursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

func NewPhonesWA(config Responser) *PhonesWA {
	if v, ok := config.(*PhonesWA); ok {
		v.ResponseType = ResponsePhonesWA
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *PhonesWA")
}
