package types

type MessageDocument struct {
	Messager `json:"messager,omitempty"`
	Header
	Media `json:"document"`
}

func NewMessageDocument(m *MessageDocument) Messager {
	mk := &messagerKernel{
		Type:             MessageTypeDocument,
		m:                m,
		Link:             m.Link,
		MessagingProduct: m.MessagingProduct,
	}

	m.Messager = mk
	return m
}
