package response

import "fmt"

type Error struct {
	KernelResponser
	Type           string `json:"type,omitempty"`
	ResponseType   string `json:"response_type,omitempty"`
	Message        string `json:"message,omitempty"`
	Code           int64  `json:"code,omitempty"`
	ErrorSubcode   int64  `json:"error_subcode,omitempty"`
	IsTransient    bool   `json:"is_transient,omitempty"`
	ErrorUserTitle string `json:"error_user_title,omitempty"`
	ErrorUserMsg   string `json:"error_user_msg,omitempty"`

	FbtraceID string `json:"fbtrace_id,omitempty"`
}

func NewError(config Responser) *Error {
	if v, ok := config.(*Error); ok {
		v.ResponseType = ResponseError
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *Error")
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error is: %s, type: %s, code: %d, ErrorSubcode: %d, FbtraceID: %s", e.Message, e.Type, e.Code, e.ErrorSubcode, e.FbtraceID)
}

// func (e *Error) UnknowResponse() string {
// 	return fmt.Sprintf("Error is: %s, type: %s, code: %d, ErrorSubcode: %d, FbtraceID: %s", e.Message, e.Type, e.Code, e.ErrorSubcode, e.FbtraceID)
// }
