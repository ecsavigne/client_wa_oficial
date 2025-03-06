package types

type MessageInteractive struct {
	Messager `json:"messager,omitempty"`
	Header
	InteractiveProto `json:"interactive"`
}

func NewMessageInteractive(m MessageInteractive) Messager {
	mk := &messagerKernel{
		Type: MessageTypeInteractive,
		m:    m,
	}

	m.Messager = mk
	return &m
}
