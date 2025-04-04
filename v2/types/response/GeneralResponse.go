package response

type GeneralResponse struct {
	KernelResponser `json:",omitempty"`
	Type            string `json:"type,omitempty"`
	*PhonesWA       `json:",omitempty"`
	*Error          `json:",omitempty"`
	*Success        `json:",omitempty"`
	*MediaInfo      `json:",omitempty"`
}

func (a *GeneralResponse) GetResponseType() ResponserRequest {
	switch {
	case a.PhonesWA != nil:
		a.PhonesWA.Type = ResponsePhonesWA
		return a
	case a.Error != nil:
		a.Error.Type = ResponseError
		return a.Error
	case a.Success != nil:
		a.Success.Type = ResponseSuccess
		return a.Success
	case a.MediaInfo != nil:
		a.MediaInfo.Type = ResponseMediaInfo
		return a.MediaInfo
	}

	return nil
}
