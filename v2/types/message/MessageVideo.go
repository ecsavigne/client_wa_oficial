package message

type MessageVideo struct {
	MessagerKernel
	Media `json:"video"`
}

func (*MessageVideo) NewVideoMessage(config Messager) *MessageVideo {
	switch v := any(config).(type) {
	case *MessageVideo:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageVideo")
}
