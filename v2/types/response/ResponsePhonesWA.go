package response

/*
Represent the response from the GET "https://graph.facebook.com/v22.0/{whatsapp-business-account-id}/phone_numbers
*/

type PhonesWA struct {
	KernelResponser
	Type         string  `json:"type,omitempty"`
	ResponseType string  `json:"response_type,omitempty"`
	Data         []Datum `json:"data"`
	Paging       Paging  `json:"paging"`
}

type Datum struct {
	VerifiedName           string               `json:"verified_name,omitempty"`
	CodeVerificationStatus string               `json:"code_verification_status,omitempty"`
	DisplayPhoneNumber     string               `json:"display_phone_number,omitempty"`
	QualityRating          string               `json:"quality_rating,omitempty"`
	PlatformType           string               `json:"platform_type,omitempty"`
	Throughput             Throughput           `json:"throughput,omitempty"`
	WebhookConfiguration   WebhookConfiguration `json:"webhook_configuration,omitempty"`
	ID                     string               `json:"id,omitempty"`
}

type Throughput struct {
	Level string `json:"level,omitempty"`
}

type WebhookConfiguration struct {
	Application string `json:"application,omitempty"`
}

type Paging struct {
	Cursors Cursors `json:"cursors,omitempty"`
}

type Cursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}
