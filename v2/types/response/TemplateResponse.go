package response

import "github.com/ecsavigne/client_wa_oficial/v2/types"

type TemplateResponse struct {
	KernelResponser
	ResponseType string                 `json:"response_type,omitempty"`
	Data         []types.MockupTemplate `json:"data"`
	Paging       *Paging                `json:"paging,omitempty"`
}

func NewTemplateResponse(config Responser) *TemplateResponse {
	if v, ok := config.(*TemplateResponse); ok {
		v.ResponseType = ResponseTemplate
		v.KernelResponser.parent = v
		return v
	}
	panic("type Responser is not *TemplateResponse")
}
