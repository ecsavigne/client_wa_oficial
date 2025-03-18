package types

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
)

type ResponseType = string

const (
	ResponseSuccess   ResponseType = "response_success"
	ResponseError     ResponseType = "response_error"
	ResponseMediaInfo ResponseType = "response_media_info"
)

type ResponserRequest interface {
	GetType() string
	String() string
	GetResponseError() *Error
	GetResponseSuccess() *Success
	GetResponseMediaInfo() *MediaInfo
}

type Message struct {
	ID            string `json:"id"`
	MessageStatus string `json:"message_status,omitempty"`
}

type Error struct {
	Type         string `json:"type,omitempty"`
	Message      string `json:"message,omitempty"`
	Code         int64  `json:"code,omitempty"`
	ErrorSubcode int64  `json:"error_subcode,omitempty"`
	FbtraceID    string `json:"fbtrace_id,omitempty"`
}

func (*Error) GetType() string { return ResponseError }

func (e *Error) String() string {
	return val(e)
}

func (e *Error) GetResponseError() *Error {
	return e
}
func (e *Error) GetResponseSuccess() *Success {
	return nil
}
func (e *Error) GetResponseMediaInfo() *MediaInfo {
	return nil
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error is: %s, type: %s, code: %d, ErrorSubcode: %d, FbtraceID: %s", e.Message, e.Type, e.Code, e.ErrorSubcode, e.FbtraceID)
}

type ContactResponse struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type Success struct {
	Type             string            `json:"type,omitempty"`
	MessagingProduct string            `json:"messaging_product,omitempty"`
	Contacts         []ContactResponse `json:"contacts,omitempty"`
	Messages         []Message         `json:"messages,omitempty"`
}

func (*Success) GetType() string { return ResponseSuccess }

func (s *Success) String() string {
	return val(s)
}

func (s *Success) GetResponseError() *Error {
	return nil
}
func (s *Success) GetResponseSuccess() *Success {
	return s
}
func (s *Success) GetResponseMediaInfo() *MediaInfo {
	return nil
}

// type ResponseRequest struct {
// 	Error
// 	Success
// 	MediaInfo
// }

func val(r any) string {
	by, msg_e := json.MarshalIndent(r, "", "  ")
	if msg_e != nil {
		panic(fmt.Sprintf("Error occurred in MarshalIndent. Error is: %s", msg_e))
	}
	return string(by)
}

type MediaInfo struct {
	Type             string `json:"type,omitempty"`
	MessagingProduct string `json:"messaging_product,omitempty"`
	MimeType         string `json:"mimetype,omitempty"`
	Sha256           string `json:"sha256,omitempty"`
	FileSize         string `json:"file_size,omitempty"`
	ID               string `json:"id,omitempty"`
	Url              string `json:"url,omitempty"`
}

func (*MediaInfo) GetType() string { return ResponseMediaInfo }

func (s *MediaInfo) String() string {
	return val(s)
}

func (m *MediaInfo) GetId() string {
	return m.ID
}

func (mi *MediaInfo) GetResponseError() *Error {
	return nil
}
func (mi *MediaInfo) GetResponseSuccess() *Success {
	return nil
}
func (mi *MediaInfo) GetResponseMediaInfo() *MediaInfo {
	return mi
}

func JsonWrapperResponseRequest(data []byte) ResponserRequest {
	// wrapper := &ResponseRequest{}
	wrapper := map[string]any{}
	err := json.Unmarshal(data, &wrapper)
	if err != nil {
		return &Error{
			Code:    401,
			Message: err.Error(),
		}
	}

	sliceKey := slices.Collect(maps.Keys(wrapper))
	isMediaInfo := !slices.ContainsFunc(sliceKey, func(v string) bool {
		mediaInfo := []string{"messaging_product", "mimetype", "sha256", "file_size", "id", "url"}
		return !slices.Contains(mediaInfo, v)
	})

	isError := !slices.ContainsFunc(sliceKey, func(v string) bool {
		errorResponse := []string{"message", "code", "error_subcode", "fbtrace_id"}
		return !slices.Contains(errorResponse, v)
	})

	isSuccess := !slices.ContainsFunc(sliceKey, func(v string) bool {
		successResponse := []string{"messaging_product", "contacts", "messages"}
		return !slices.Contains(successResponse, v)
	})

	switch {
	case isMediaInfo:
		mInfo := &MediaInfo{}
		b, _ := json.Marshal(wrapper)
		json.Unmarshal(b, mInfo)
		return mInfo
	case isError:
		errorResponse := &Error{}
		b, _ := json.Marshal(wrapper)
		json.Unmarshal(b, errorResponse)
		return errorResponse
	case isSuccess:
		successResponse := &Success{}
		b, _ := json.Marshal(wrapper)
		json.Unmarshal(b, successResponse)
		return successResponse
	}

	if _, ok := wrapper["error"]; ok {
		return &Error{
			Message: wrapper["error"].(string),
			Code:    401,
		}
	}
	return nil
}
