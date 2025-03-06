package types

type MessageAudio struct {
	Messager `json:"messager,omitempty"`
	Header
	Media `json:"audio"`
}

func NewMessageAudio(m MessageAudio) Messager {
	mk := &messagerKernel{
		Type: MessageTypeAudio,
		m:    m,
	}

	m.Messager = mk
	return &m
}
