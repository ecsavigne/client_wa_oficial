package types

type Text struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type MessageText struct {
	Messager `json:"messager,omitempty"`
	Header
	Text `json:"text"`
}

func NewMessageText(m *MessageText) Messager {
	mk := &messagerKernel{
		Type: MessageTypeText,
		m:    m,
	}

	m.Messager = mk
	return m
}
