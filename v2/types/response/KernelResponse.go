package response

type ResponseType = string

const (
	ResponseSuccess         ResponseType = "response_success"
	ResponseError           ResponseType = "response_error"
	ResponseMediaInfo       ResponseType = "response_media_info"
	ResponseGeneralResponse ResponseType = "response_general_response"
	ResponsePhonesWA        ResponseType = "response_phones_wa"
	ResponsePhone           ResponseType = "response_phone"
	ResponseWABA            ResponseType = "response_waba"
	ResponseBusiness        ResponseType = "response_business"
	ResponseTemplate        ResponseType = "response_template"
	ResponseMockupTemplate  ResponseType = "response_mockup_template"
	ResponseWebHookTemplate ResponseType = "response_webhook_template"
	ResponseUnknow          ResponseType = "response_unknow"
	ResponseOther           ResponseType = "response_other"
)

type Responser interface {
	GetType() string
	String() string
	GetResponseError() *Error
	GetResponseSuccess() *Success
	GetResponseMediaInfo() *MediaInfo
	GetResponsePhonesWA() *PhonesWA
	GetResponsePhone() *Phone
	GetResponseWaba() *Waba
	GetResponseBusiness() *Business
	GetGeneralResponse() *GeneralResponse
	GetTemplateResponse() *TemplateResponse
	GetMockupTemplateResponse() *MockupTemplateResponse
	GetWebHookTemplateResponse() *WebHookTemplateResponse
	GetUnknowResponse() *UnknowResponse
	IsType(ResponseType) bool
}

type Message struct {
	ID            string `json:"id"`
	MessageStatus string `json:"message_status,omitempty"`
}

type ContactResponse struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type KernelResponser struct {
	parent Responser
}

// GetType returns the type of the ResponserRequest interface.
// It returns a string with the type of the interface.
// If the interface is not of type *Error, *Success, *MediaInfo, *PhonesWA or *GeneralResponse,
// it returns an empty string.
func (k *KernelResponser) GetType() string {
	switch v := k.parent.(type) {
	case *Error:
		return v.ResponseType
	case *Success:
		return v.ResponseType
	case *MediaInfo:
		return v.ResponseType
	case *PhonesWA:
		return v.ResponseType
	case *Phone:
		return v.ResponseType
	case *Waba:
		return v.ResponseType
	case *GeneralResponse:
		return v.ResponseType
	case *TemplateResponse:
		return v.ResponseType
	case *MockupTemplateResponse:
		return v.ResponseType
	case *WebHookTemplateResponse:
		return v.ResponseType
	case *UnknowResponse:
		return v.ResponseType
	default:
		return ""
	}
}

// String returns a JSON-formatted string representation of the KernelResponser
// object. It determines the specific type of KernelResponser and uses the val
// function to convert it to a string. If the type is not recognized, it returns
// an empty string.

func (k *KernelResponser) String() string {
	switch v := k.parent.(type) {
	case *Error:
		return Val(v)
	case *Success:
		return Val(v)
	case *MediaInfo:
		return Val(v)
	case *PhonesWA:
		return Val(v)
	case *Phone:
		return Val(v)
	case *Waba:
		return Val(v)
	case *Business:
		return Val(v)
	case *GeneralResponse:
		return Val(v)
	case *TemplateResponse:
		return Val(v)
	case *MockupTemplateResponse:
		return Val(v)
	case *WebHookTemplateResponse:
		return Val(v)
	case *UnknowResponse:
		return Val(v)
	default:
		return ""
	}
}

// GetResponseError attempts to cast the KernelResponser to an *Error type.
// If successful, it returns the *Error instance; otherwise, it returns nil.

func (k *KernelResponser) GetResponseError() *Error {
	if v, ok := k.parent.(*Error); ok {
		return v
	}
	return nil
}

// GetResponseSuccess attempts to cast the KernelResponser to an *Success type.
// If successful, it returns the *Success instance; otherwise, it returns nil.
func (k *KernelResponser) GetResponseSuccess() *Success {
	if v, ok := k.parent.(*Success); ok {
		return v
	}
	return nil
}

// GetResponseMediaInfo attempts to cast the KernelResponser to an *MediaInfo type.
// If successful, it returns the *MediaInfo instance; otherwise, it returns nil.
func (k *KernelResponser) GetResponseMediaInfo() *MediaInfo {
	if v, ok := k.parent.(*MediaInfo); ok {
		return v
	}
	return nil
}

// GetResponsePhonesWA attempts to cast the KernelResponser to a *PhonesWA type.
// If successful, it returns the *PhonesWA instance; otherwise, it returns nil.

func (k *KernelResponser) GetResponsePhonesWA() *PhonesWA {
	if v, ok := k.parent.(*PhonesWA); ok {
		return v
	}
	return nil
}

func (k *KernelResponser) GetResponsePhone() *Phone {
	if v, ok := k.parent.(*Phone); ok {
		return v
	}
	return nil
}

func (k *KernelResponser) GetResponseWaba() *Waba {
	if v, ok := k.parent.(*Waba); ok {
		return v
	}
	return nil
}

func (k *KernelResponser) GetResponseBusiness() *Business {
	if v, ok := k.parent.(*Business); ok {
		return v
	}
	return nil
}

// GetGeneralResponse attempts to cast the KernelResponser to a *GeneralResponse type.
// If successful, it returns the *GeneralResponse instance; otherwise, it returns nil.
func (k *KernelResponser) GetGeneralResponse() *GeneralResponse {
	if v, ok := k.parent.(*GeneralResponse); ok {
		return v
	}
	return nil
}

// GetGeneralResponse attempts to cast the KernelResponser to a *GeneralResponse type.
// If successful, it returns the *GeneralResponse instance; otherwise, it returns nil.
func (k *KernelResponser) GetTemplateResponse() *TemplateResponse {
	if v, ok := k.parent.(*TemplateResponse); ok {
		return v
	}
	return nil
}

// GetMockupTemplateResponse attempts to cast the KernelResponser to a *MockupTemplateResponse type.
// If successful, it returns the *MockupTemplateResponse instance; otherwise, it returns nil.
func (k *KernelResponser) GetMockupTemplateResponse() *MockupTemplateResponse {
	if v, ok := k.parent.(*MockupTemplateResponse); ok {
		return v
	}
	return nil
}

// GetWebHookTemplateResponse attempts to cast the KernelResponser to a *WebHookTemplateResponse type.
// If successful, it returns the *WebHookTemplateResponse instance; otherwise, it returns nil.
func (k *KernelResponser) GetWebHookTemplateResponse() *WebHookTemplateResponse {
	if v, ok := k.parent.(*WebHookTemplateResponse); ok {
		return v
	}
	return nil
}

// GetUnknowResponse attempts to cast the KernelResponser to a *UnknowResponse type.
// If successful, it returns the *UnknowResponse instance; otherwise, it returns nil.
func (k *KernelResponser) GetUnknowResponse() *UnknowResponse {
	if v, ok := k.parent.(*UnknowResponse); ok {
		return v
	}
	return nil
}

// IsType returns true if the KernelResponser type matches the given ResponseType, and false otherwise.
func (k *KernelResponser) IsType(pType ResponseType) bool {
	return k.GetType() == pType
}
