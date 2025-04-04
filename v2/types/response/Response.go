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
)

type ResponserRequest interface {
	GetType() string
	String() string
	GetResponseError() *Error
	GetResponseSuccess() *Success
	GetResponseMediaInfo() *MediaInfo
	GetResponsePhonesWA() *PhonesWA
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
func JsonWrapperResponseRequest(data []byte) ResponserRequest {
	wrapper := map[string]any{}
	err := json.Unmarshal(data, &wrapper)
	if err != nil {
		return &Error{
			Code:    401,
			Message: err.Error(),
		}
	}

	if _, ok := wrapper["error"]; ok {
		return &Error{
			Message: wrapper["error"].(string),
			Code:    401,
		}
	}

	gralResponse := &GeneralResponse{}
	b, _ := json.Marshal(wrapper)
	json.Unmarshal(b, gralResponse)
	return gralResponse.GetResponseType()
}

type KernelResponser struct{}

// GetType returns the type of the ResponserRequest interface.
// It returns a string with the type of the interface.
// If the interface is not of type *Error, *Success, *MediaInfo, *PhonesWA or *GeneralResponse,
// it returns an empty string.
func (k *KernelResponser) GetType() string {
	switch v := any(k).(type) {
	case *Error:
		return v.Type
	case *Success:
		return v.Type
	case *MediaInfo:
		return v.Type
	case *PhonesWA:
		return v.Type
	case *GeneralResponse:
		return v.Type
	default:
		return ""
	}
}

// String returns a JSON-formatted string representation of the KernelResponser
// object. It determines the specific type of KernelResponser and uses the val
// function to convert it to a string. If the type is not recognized, it returns
// an empty string.

func (k *KernelResponser) String() string {
	switch v := any(k).(type) {
	case *Error:
		return Val(v)
	case *Success:
		return Val(v)
	case *MediaInfo:
		return Val(v)
	case *PhonesWA:
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
	if v, ok := any(k).(*Error); ok {
		return v
	}
	return nil
}

// GetResponseSuccess attempts to cast the KernelResponser to an *Success type.
// If successful, it returns the *Success instance; otherwise, it returns nil.
func (k *KernelResponser) GetResponseSuccess() *Success {
	if v, ok := any(k).(*Success); ok {
		return v
	}
	return nil
}

// GetResponseMediaInfo attempts to cast the KernelResponser to an *MediaInfo type.
// If successful, it returns the *MediaInfo instance; otherwise, it returns nil.
func (k *KernelResponser) GetResponseMediaInfo() *MediaInfo {
	if v, ok := any(k).(*MediaInfo); ok {
		return v
	}
	return nil
}

// GetResponsePhonesWA attempts to cast the KernelResponser to a *PhonesWA type.
// If successful, it returns the *PhonesWA instance; otherwise, it returns nil.

func (k *KernelResponser) GetResponsePhonesWA() *PhonesWA {
	if v, ok := any(k).(*PhonesWA); ok {
		return v
	}
	return nil
}

// GetGeneralResponse attempts to cast the KernelResponser to a *GeneralResponse type.
// If successful, it returns the *GeneralResponse instance; otherwise, it returns nil.
func (k *KernelResponser) GetGeneralResponse() *GeneralResponse {
	if v, ok := any(k).(*GeneralResponse); ok {
		return v
	}
	return nil
}

// IsType returns true if the KernelResponser type matches the given ResponseType, and false otherwise.
func (k *KernelResponser) IsType(pType ResponseType) bool {
	return k.GetType() == pType
}

type Error struct {
	KernelResponser `json:",omitempty"`
	Type            string `json:"type,omitempty"`
	Message         string `json:"message,omitempty"`
	Code            int64  `json:"code,omitempty"`
	ErrorSubcode    int64  `json:"error_subcode,omitempty"`
	FbtraceID       string `json:"fbtrace_id,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error is: %s, type: %s, code: %d, ErrorSubcode: %d, FbtraceID: %s", e.Message, e.Type, e.Code, e.ErrorSubcode, e.FbtraceID)
}

type Success struct {
	KernelResponser  `json:",omitempty"`
	Type             string            `json:"type,omitempty"`
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
	KernelResponser  `json:",omitempty"`
	Type             string `json:"type,omitempty"`
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
