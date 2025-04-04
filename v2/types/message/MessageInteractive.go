package message

type MessageInteractive struct {
	MessagerKernel
	InteractiveProto `json:"interactive"`
}

// IsType returns true if the MessageInteractive type matches the given type, and false otherwise
// type is:  button, list,product,catalog_message, product_list,flow,cta_url
func (i *MessageInteractive) IsType(t string) bool {
	return i.InteractiveProto.Type == t
}
