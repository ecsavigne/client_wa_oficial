package response

type Phone struct {
	KernelResponser
	ResponseType string `json:"response_type,omitempty"`
	*PhoneInfo
}

func NewPhone(config ResponserRequest) *Phone {
	if v, ok := config.(*Phone); ok {
		v.ResponseType = ResponsePhone
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *Phone")
}
