package response

import "fmt"

type UnknowResponse struct {
	KernelResponser
	ResponseType string         `json:"response_type,omitempty"`
	Unknow       map[string]any `json:"unknow,omitempty"`
}

func NewUnknowResponse(config Responser) *UnknowResponse {
	if v, ok := config.(*UnknowResponse); ok {
		v.ResponseType = ResponseUnknow
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *Error")
}

func (e *Error) UnknowResponse() string {
	return fmt.Sprintf("Error is: %s, type: %s, code: %d, ErrorSubcode: %d, FbtraceID: %s", e.Message, e.Type, e.Code, e.ErrorSubcode, e.FbtraceID)
}
