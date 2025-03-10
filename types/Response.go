package types

import (
	"encoding/json"
	"fmt"
)

const (
	ResponseSuccess = "response_success"
	ResponseError   = "response_error"
)

type ResponserRequest interface {
	GetType() string
	String() string
}

type Message struct {
	ID            string `json:"id"`
	MessageStatus string `json:"message_status,omitempty"`
}

type Error struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int64  `json:"code"`
	ErrorSubcode int64  `json:"error_subcode"`
	FbtraceID    string `json:"fbtrace_id"`
}

func (*Error) GetType() string { return ResponseError }

func (e *Error) String() string {
	return val(e)
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error is: %s, type: %s, code: %d, ErrorSubcode: %d, FbtraceID: %s", e.Message, e.Type, e.Code, e.ErrorSubcode, e.FbtraceID)
}

type ContactResponse struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type Success struct {
	MessagingProduct string            `json:"messaging_product"`
	Contacts         []ContactResponse `json:"contacts"`
	Messages         []Message         `json:"messages"`
	Id               string            `json:"id,omitempty"`
}

func (*Success) GetType() string { return ResponseSuccess }

func (s *Success) String() string {
	return val(s)
}

type ResponseRequest struct {
	*Error `json:"error,omitempty"`
	*Success
}

func val(r any) string {
	by, msg_e := json.MarshalIndent(r, "", "  ")
	if msg_e != nil {
		panic(fmt.Sprintf("Error occurred in MarshalIndent. Error is: %s", msg_e))
	}
	return string(by)
}
