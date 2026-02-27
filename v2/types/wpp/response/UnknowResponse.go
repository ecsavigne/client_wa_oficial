package response

type UnknowResponse struct {
	KernelResponser
	ResponseType string         `json:"response_type,omitempty"`
	Unknow       map[string]any `json:"unknow,omitempty"`
}

func NewUnknowResponse(config Responser) *UnknowResponse {
	if v, ok := config.(*UnknowResponse); ok {
		v.ResponseType = ResponseUnknow
		if v.Unknow == nil {
			v.Unknow = make(map[string]any)
		}
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *UnknowResponse")
}
