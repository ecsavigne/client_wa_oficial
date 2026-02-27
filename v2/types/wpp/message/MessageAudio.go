package message

type MessageAudio struct {
	MessagerKernel
	Media `json:"audio"`
}
