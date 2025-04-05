package message

type MessageDocument struct {
	MessagerKernel
	Media `json:"document"`
}

func (*MessageDocument) NewDocumentMessage(config Messager) *MessageDocument {
	switch v := any(config).(type) {
	case *MessageDocument:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageDocument")
}
