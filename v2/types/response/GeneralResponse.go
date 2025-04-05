package response

type GeneralResponse struct {
	KernelResponser
	Type       string `json:"type,omitempty"`
	*PhonesWA  `json:",omitempty"`
	*Error     `json:",omitempty"`
	*Success   `json:",omitempty"`
	*MediaInfo `json:",omitempty"`
}

func (a *GeneralResponse) GetResponseType() ResponserRequest {
	switch {
	case a.PhonesWA != nil:
		a.PhonesWA.Type = ResponsePhonesWA
		return NewPhonesWA(a.PhonesWA)
	case a.Error != nil:
		a.Error.Type = ResponseError
		return NewError(a.Error)
	case a.Success != nil:
		a.Success.Type = ResponseSuccess
		return NewSuccess(a.Success)
	case a.MediaInfo != nil:
		a.MediaInfo.Type = ResponseMediaInfo
		return NewMediaInfo(a.MediaInfo)
	}

	return nil
}
