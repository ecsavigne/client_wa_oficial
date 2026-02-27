package response

import "github.com/ecsavigne/client_wa_oficial/v2/types/wpp"

type MockupTemplateResponse struct {
	KernelResponser
	ResponseType string               `json:"response_type,omitempty"`
	Data         []wpp.MockupTemplate `json:"data"`
	Paging       *Paging              `json:"paging,omitempty"`
}

func NewMockupTemplateResponse(config Responser) *MockupTemplateResponse {
	if v, ok := config.(*MockupTemplateResponse); ok {
		v.ResponseType = ResponseTemplate
		v.KernelResponser.parent = v
		return v
	}
	panic("type Responser is not *MockupTemplateResponse")
}
