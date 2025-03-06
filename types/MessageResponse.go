package types

type Context struct {
	MessageId string `json:"message_id" valid:"required"`
}

type MessageResponse struct {
	Messager `json:"messager,omitempty"`
	Header
	Context `json:"context"`
	Text    `json:"text"`
}

func NewMessageResponse(m MessageResponse) Messager {
	mk := &messagerKernel{
		Type: "response",
		m:    m,
	}

	m.Messager = mk
	return &m
}
