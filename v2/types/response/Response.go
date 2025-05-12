package response

import (
	"encoding/json"
	"fmt"
)

type ResponseType = string

const (
	ResponseSuccess         ResponseType = "response_success"
	ResponseError           ResponseType = "response_error"
	ResponseMediaInfo       ResponseType = "response_media_info"
	ResponseGeneralResponse ResponseType = "response_general_response"
	ResponsePhonesWA        ResponseType = "response_phones_wa"
	ResponseWABA            ResponseType = "response_waba"
)

type ResponserRequest interface {
	GetType() string
	String() string
	GetResponseError() *Error
	GetResponseSuccess() *Success
	GetResponseMediaInfo() *MediaInfo
	GetResponsePhonesWA() *PhonesWA
	GetResponseWaba() *Waba
	GetGeneralResponse() *GeneralResponse
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

func NewError(config ResponserRequest) *Error {
	if v, ok := config.(*Error); ok {
		v.ResponseType = ResponseError
		v.KernelResponser.parent = v
		return v
	}
	return nil
}

func NewSuccess(config ResponserRequest) *Success {
	if v, ok := config.(*Success); ok {
		v.ResponseType = ResponseSuccess
		v.KernelResponser.parent = v
		return v
	}
	return nil
}

func NewMediaInfo(config ResponserRequest) *MediaInfo {
	if v, ok := config.(*MediaInfo); ok {
		v.ResponseType = ResponseMediaInfo
		v.KernelResponser.parent = v
		return v
	}
	return nil
}

func NewPhonesWA(config ResponserRequest) *PhonesWA {
	if v, ok := config.(*PhonesWA); ok {
		v.ResponseType = ResponsePhonesWA
		v.KernelResponser.parent = v
		return v
	}
	return nil
}

func NewWABA(config ResponserRequest) *Waba {
	if v, ok := config.(*Waba); ok {
		v.ResponseType = ResponseWABA
		v.KernelResponser.parent = v
		return v
	}
	return nil
}

func NewGeneralResponse(config ResponserRequest) *GeneralResponse {
	if v, ok := config.(*GeneralResponse); ok {
		v.ResponseType = ResponseGeneralResponse
		v.KernelResponser.parent = v
		return v
	}
	return nil
}

func Val(r any) string {
	by, msg_e := json.MarshalIndent(r, "", "  ")
	if msg_e != nil {
		panic(fmt.Sprintf("Error occurred in MarshalIndent. Error is: %s", msg_e))
	}
	return string(by)
}

// JsonWrapperResponseRequest is a function that wraps a given json data into a ResponseWrapper.
// It unmarshals the data into a map[string]any, and then checks if the map contains an "error" key.
// If it does, it returns a types.Error object with the error message and code 401.
// Otherwise, it marshals the map back into json and unmarshals it into a types.GeneralResponse object.
// Finally, return interface que represent of type of response.
func JsonWrapperResponseRequest(dataBin []byte) ResponserRequest {
	wrapper := map[string]any{}
	gralResponse := NewGeneralResponse(&GeneralResponse{})

	err := json.Unmarshal(dataBin, &wrapper)
	if err != nil {
		return NewError(&Error{
			Code:    401,
			Message: err.Error(),
		})
	}

	if _, ok := wrapper["error"]; ok {
		switch wrapper["error"].(type) {
		case string:
			return NewError(&Error{
				Message: wrapper["error"].(string),
				Code:    401,
			})
		default:
			data, _ := json.Marshal(wrapper["error"])
			errorResponse := NewError(&Error{})
			json.Unmarshal(data, errorResponse)
			gralResponse.Error = errorResponse
			return gralResponse.GetResponseType()
		}

	}

	data, _ := json.Marshal(wrapper)

	if wrapper["data"] != nil && wrapper["paging"] != nil {
		phonesWA := NewPhonesWA(&PhonesWA{})
		json.Unmarshal(data, phonesWA)
		if len(phonesWA.Data) > 0 && (phonesWA.Data[0].DisplayPhoneNumber != "" && phonesWA.Data[0].VerifiedName != "") {
			gralResponse.PhonesWA = phonesWA
			return gralResponse.GetResponseType()
		} else {
			waba := NewWABA(&Waba{})
			json.Unmarshal(data, waba)
			if waba != nil {
				gralResponse.Waba = waba
				return gralResponse.GetResponseType()
			}
		}
	}

	json.Unmarshal(data, gralResponse)
	return gralResponse.GetResponseType()
}

type KernelResponser struct {
	parent ResponserRequest
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
	case *Waba:
		return v.ResponseType
	case *GeneralResponse:
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
	case *Waba:
		return Val(v)
	case *GeneralResponse:
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

func (k *KernelResponser) GetResponseWaba() *Waba {
	if v, ok := k.parent.(*Waba); ok {
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

// IsType returns true if the KernelResponser type matches the given ResponseType, and false otherwise.
func (k *KernelResponser) IsType(pType ResponseType) bool {
	return k.GetType() == pType
}

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

func (e *Error) Error() string {
	return fmt.Sprintf("Error is: %s, type: %s, code: %d, ErrorSubcode: %d, FbtraceID: %s", e.Message, e.Type, e.Code, e.ErrorSubcode, e.FbtraceID)
}

type Success struct {
	KernelResponser
	Type             string            `json:"type,omitempty"`
	ResponseType     string            `json:"response_type,omitempty"`
	MessagingProduct string            `json:"messaging_product,omitempty"`
	Contacts         []ContactResponse `json:"contacts,omitempty"`
	Messages         []Message         `json:"messages,omitempty"`
	Success          bool              `json:"success,omitempty"`
	MediaInfo        *MediaInfo        `json:"media_info,omitempty"`
}

func (s *Success) GetMediaInfo() *MediaInfo {
	return s.MediaInfo
}

func (s *Success) GetMessageId() string {
	return s.Messages[0].ID
}

type MediaInfo struct {
	KernelResponser
	Type             string `json:"type,omitempty"`
	ResponseType     string `json:"response_type,omitempty"`
	MessagingProduct string `json:"messaging_product,omitempty"`
	MimeType         string `json:"mime_type,omitempty"`
	Sha256           string `json:"sha256,omitempty"`
	FileSize         uint64 `json:"file_size,omitempty"`
	ID               string `json:"id,omitempty"`
	Url              string `json:"url,omitempty"`
}

func (m *MediaInfo) GetId() string {
	return m.ID
}
