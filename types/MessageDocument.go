package types

type MessageDocument struct {
	Messager `json:"messager,omitempty"`
	Header
	Media `json:"document"`
}

func NewMessageDocument(m MessageDocument) Messager {
	mk := &messagerKernel{
		Type: MessageTypeDocument,
		m:    m,
	}

	m.Messager = mk
	return &m
}
