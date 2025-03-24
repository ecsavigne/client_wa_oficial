package types

/*
Represent the response from the GET "https://graph.facebook.com/v22.0/{whatsapp-business-account-id}/phone_numbers
*/

type PhonesWA struct {
	Data   []Datum `json:"data"`
	Paging Paging  `json:"paging"`
}

type Datum struct {
	VerifiedName           string               `json:"verified_name"`
	CodeVerificationStatus string               `json:"code_verification_status"`
	DisplayPhoneNumber     string               `json:"display_phone_number"`
	QualityRating          string               `json:"quality_rating"`
	PlatformType           string               `json:"platform_type"`
	Throughput             Throughput           `json:"throughput"`
	WebhookConfiguration   WebhookConfiguration `json:"webhook_configuration"`
	ID                     string               `json:"id"`
}

type Throughput struct {
	Level string `json:"level"`
}

type WebhookConfiguration struct {
	Application string `json:"application"`
}

type Paging struct {
	Cursors Cursors `json:"cursors"`
}

type Cursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

func (a *PhonesWA) GetType() string {
	return ResponsePhonesWA
}
func (a *PhonesWA) String() string {
	return val(a)
}

func (a *PhonesWA) GetResponsePhonesWA() *PhonesWA {
	return a
}

func (a *PhonesWA) GetGeneralResponse() *GeneralResponse {
	return nil
}

func (a *PhonesWA) GetResponseError() *Error {
	return nil
}

func (a *PhonesWA) GetResponseSuccess() *Success {
	return nil
}

func (a *PhonesWA) GetResponseMediaInfo() *MediaInfo {
	return nil
}
