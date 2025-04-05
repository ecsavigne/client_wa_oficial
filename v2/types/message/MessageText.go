package message

type Text struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type MessageText struct {
	MessagerKernel
	Text `json:"text"`
}
