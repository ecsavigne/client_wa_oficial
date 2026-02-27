package response

type OtherResponse struct {
	KernelResponser
	ResponseType string         `json:"response_type,omitempty"`
	Other        map[string]any `json:",omitempty"`
}

func NewOtherResponse(config Responser) *OtherResponse {
	if v, ok := config.(*OtherResponse); ok {
		v.ResponseType = ResponseOther
		if v.Other == nil {
			v.Other = make(map[string]any)
		}
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *Other")
}
