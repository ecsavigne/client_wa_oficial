package response

type Waba struct {
	KernelResponser
	Type         string     `json:"type,omitempty"`
	ResponseType string     `json:"response_type,omitempty"`
	Data         []WabaInfo `json:"data"`
	Paging       Paging     `json:"paging"`
}

type WabaInfo struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	Currency                 string `json:"currency"`
	TimezoneID               string `json:"timezone_id"`
	MessageTemplateNamespace string `json:"message_template_namespace"`
}

func NewWABA(config ResponserRequest) *Waba {
	if v, ok := config.(*Waba); ok {
		v.ResponseType = ResponseWABA
		v.KernelResponser.parent = v
		return v
	}
	return nil
}
