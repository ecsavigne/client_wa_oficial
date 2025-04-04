package message

type MessageDocument struct {
	MessagerKernel
	Media `json:"document"`
}
