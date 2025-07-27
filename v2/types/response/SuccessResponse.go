package response

type Success struct {
	KernelResponser
	Type             string            `json:"type,omitempty"`
	ResponseType     string            `json:"response_type,omitempty"`
	MessagingProduct string            `json:"messaging_product,omitempty"`
	Contacts         []ContactResponse `json:"contacts,omitempty"`
	Messages         []Message         `json:"messages,omitempty"`
	Success          bool              `json:"success,omitempty"`
	MediaInfo        *MediaInfo        `json:"media_info,omitempty"`
}

func NewSuccess(config Responser) *Success {
	if v, ok := config.(*Success); ok {
		v.ResponseType = ResponseSuccess
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *Success")
}

func (s *Success) GetMediaInfo() *MediaInfo {
	return s.MediaInfo
}

func (s *Success) GetMessageId() string {
	return s.Messages[0].ID
}
