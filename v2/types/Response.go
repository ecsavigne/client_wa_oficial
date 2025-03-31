package types

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

func val(r any) string {
	by, msg_e := json.MarshalIndent(r, "", "  ")
	if msg_e != nil {
		panic(fmt.Sprintf("Error occurred in MarshalIndent. Error is: %s", msg_e))
	}
	return string(by)
}

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

func (k *KernelResponser) String() string {
	switch v := any(k).(type) {
	case *Error:
		return val(v)
	case *Success:
		return val(v)
	case *MediaInfo:
		return val(v)
	case *PhonesWA:
		return val(v)
	case *GeneralResponse:
		return val(v)
	default:
		return ""
	}
}

func (k *KernelResponser) GetResponseError() *Error {
	if v, ok := any(k).(*Error); ok {
		return v
	}
	return nil
}

func (k *KernelResponser) GetResponseSuccess() *Success {
	if v, ok := any(k).(*Success); ok {
		return v
	}
	return nil
}

func (k *KernelResponser) GetResponseMediaInfo() *MediaInfo {
	if v, ok := any(k).(*MediaInfo); ok {
		return v
	}
	return nil
}

func (k *KernelResponser) GetResponsePhonesWA() *PhonesWA {
	if v, ok := any(k).(*PhonesWA); ok {
		return v
	}
	return nil
}

func (k *KernelResponser) GetGeneralResponse() *GeneralResponse {
	if v, ok := any(k).(*GeneralResponse); ok {
		return v
	}
	return nil
}

func (k *KernelResponser) IsType(pType ResponseType) bool {
	return k.GetType() == pType
}

type Error struct {
	ResponserRequest `json:",omitempty"`
	Type             string `json:"type,omitempty"`
	Message          string `json:"message,omitempty"`
	Code             int64  `json:"code,omitempty"`
	ErrorSubcode     int64  `json:"error_subcode,omitempty"`
	FbtraceID        string `json:"fbtrace_id,omitempty"`
}

// func (*Error) GetType() string { return ResponseError }

// func (e *Error) String() string {
// 	return val(e)
// }

// func (e *Error) GetResponseError() *Error {
// 	return e
// }
// func (e *Error) GetResponseSuccess() *Success {
// 	return nil
// }
// func (e *Error) GetResponseMediaInfo() *MediaInfo {
// 	return nil
// }

// func (e *Error) GetGeneralResponse() *GeneralResponse {
// 	return nil
// }

// func (e *Error) GetResponsePhonesWA() *PhonesWA {
// 	return nil
// }

func (e *Error) Error() string {
	return fmt.Sprintf("Error is: %s, type: %s, code: %d, ErrorSubcode: %d, FbtraceID: %s", e.Message, e.Type, e.Code, e.ErrorSubcode, e.FbtraceID)
}

// func (e *Error) IsType(t ResponseType) bool {
// 	return ResponseError == t
// }

type Success struct {
	ResponserRequest `json:",omitempty"`
	Type             string            `json:"type,omitempty"`
	MessagingProduct string            `json:"messaging_product,omitempty"`
	Contacts         []ContactResponse `json:"contacts,omitempty"`
	Messages         []Message         `json:"messages,omitempty"`
	Success          bool              `json:"success,omitempty"`
	MediaInfo        *MediaInfo        `json:"media_info,omitempty"`
}

func (*Success) GetType() string { return ResponseSuccess }

// func (s *Success) String() string {
// 	return val(s)
// }

// func (s *Success) IsType(t ResponseType) bool {
// 	return ResponseSuccess == t
// }

// func (s *Success) GetResponseError() *Error {
// 	return nil
// }
// func (s *Success) GetResponseSuccess() *Success {
// 	return s
// }
// func (s *Success) GetResponseMediaInfo() *MediaInfo {
// 	return nil
// }

// func (s *Success) GetGeneralResponse() *GeneralResponse {
// 	return nil
// }

// func (s *Success) GetResponsePhonesWA() *PhonesWA {
// 	return nil
// }

// func (s *Success) GetMediaInfo() *MediaInfo {
// 	return s.MediaInfo
// }

// func (s *Success) GetMessageId() string {
// 	return s.Messages[0].ID
// }

type MediaInfo struct {
	ResponserRequest `json:",omitempty"`
	Type             string `json:"type,omitempty"`
	MessagingProduct string `json:"messaging_product,omitempty"`
	MimeType         string `json:"mime_type,omitempty"`
	Sha256           string `json:"sha256,omitempty"`
	FileSize         uint64 `json:"file_size,omitempty"`
	ID               string `json:"id,omitempty"`
	Url              string `json:"url,omitempty"`
}

func (*MediaInfo) GetType() string { return ResponseMediaInfo }

// func (s *MediaInfo) String() string {
// 	return val(s)
// }

// func (m *MediaInfo) GetId() string {
// 	return m.ID
// }

// func (m *MediaInfo) IsType(t ResponseType) bool {
// 	return ResponseMediaInfo == t
// }

// func (mi *MediaInfo) GetResponseError() *Error {
// 	return nil
// }
// func (mi *MediaInfo) GetResponseSuccess() *Success {
// 	return nil
// }
// func (mi *MediaInfo) GetResponseMediaInfo() *MediaInfo {
// 	return mi
// }

// func (mi *MediaInfo) GetResponsePhonesWA() *PhonesWA {
// 	return nil
// }

// func (mi *MediaInfo) GetGeneralResponse() *GeneralResponse {
// 	return nil
// }
