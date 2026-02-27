package response

type WabaIdentyfy struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type WhatsappBusinessApiData struct {
	*WabaIdentyfy `json:",omitempty"`
	Link          string `json:"link"`
}

type WabaInfo struct {
	*WabaIdentyfy            `json:",omitempty"`
	Currency                 string                   `json:"currency,omitempty"`
	TimezoneID               string                   `json:"timezone_id,omitempty"`
	MessageTemplateNamespace string                   `json:"message_template_namespace,omitempty"`
	WhatsappBusinessApiData  *WhatsappBusinessApiData `json:"whatsapp_business_api_data,omitempty"`
}

type Waba struct {
	KernelResponser
	Type         string     `json:"type,omitempty"`
	ResponseType string     `json:"response_type,omitempty"`
	Data         []WabaInfo `json:"data,omitempty"`
	Paging       *Paging    `json:"paging,omitempty"`
}

func NewWABA(config Responser) *Waba {
	if v, ok := config.(*Waba); ok {
		v.ResponseType = ResponseWABA
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *Waba")
}
