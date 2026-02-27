package response

type MediaInfo struct {
	KernelResponser
	Type             string `json:"type,omitempty"`
	ResponseType     string `json:"response_type,omitempty"`
	MessagingProduct string `json:"messaging_product,omitempty"`
	MimeType         string `json:"mime_type,omitempty"`
	Sha256           string `json:"sha256,omitempty"`
	FileSize         uint64 `json:"file_size,omitempty"`
	ID               string `json:"id,omitempty"`
	Url              string `json:"url,omitempty"`
}

func NewMediaInfo(config Responser) *MediaInfo {
	if v, ok := config.(*MediaInfo); ok {
		v.ResponseType = ResponseMediaInfo
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *MediaInfo")
}

func (m *MediaInfo) GetId() string {
	return m.ID
}
