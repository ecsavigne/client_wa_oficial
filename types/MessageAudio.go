package types

type MessageAudio struct {
	Messager `json:"messager,omitempty"`
	Header
	Media `json:"audio"`
}

func NewMessageAudio(m *MessageAudio) Messager {
	mk := &messagerKernel{
		Type:             MessageTypeAudio,
		m:                m,
		Link:             m.Link,
		MessagingProduct: m.MessagingProduct,
	}

	m.Messager = mk
	return m
}
