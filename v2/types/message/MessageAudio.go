package message

type MessageAudio struct {
	MessagerKernel
	Media `json:"audio"`
}

func (*MessageAudio) NewAudioMessage(config Messager) *MessageAudio {
	switch v := any(config).(type) {
	case *MessageAudio:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageAudio")
}
