package response

type Waba struct {
	KernelResponser
	Type         string     `json:"type,omitempty"`
	ResponseType string     `json:"response_type,omitempty"`
	Data         []DataWaba `json:"data"`
	Paging       Paging     `json:"paging"`
}

type DataWaba struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	Currency                 string `json:"currency"`
	TimezoneID               string `json:"timezone_id"`
	MessageTemplateNamespace string `json:"message_template_namespace"`
}
