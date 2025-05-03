package response

type GeneralResponse struct {
	KernelResponser
	ResponseType string `json:"response_type,omitempty"`
	*PhonesWA    `json:",omitempty"`
	*Waba        `json:",omitempty"`
	*Error       `json:",omitempty"`
	*Success     `json:",omitempty"`
	*MediaInfo   `json:",omitempty"`
}

func (a *GeneralResponse) GetResponseType() ResponserRequest {
	switch {
	case a.PhonesWA != nil:
		// a.PhonesWA.Type = ResponsePhonesWA
		a.ResponseType = ResponsePhonesWA
		return NewPhonesWA(a.PhonesWA)
	case a.Waba != nil:
		// a.Waba.Type = ResponseWABA
		a.ResponseType = ResponseWABA
		return NewWABA(a.Waba)
	case a.Error != nil:
		// a.Error.Type = ResponseError
		a.ResponseType = ResponseError
		return NewError(a.Error)
	case a.Success != nil:
		// a.Success.Type = ResponseSuccess
		a.ResponseType = ResponseSuccess
		return NewSuccess(a.Success)
	case a.MediaInfo != nil:
		// a.MediaInfo.Type = ResponseMediaInfo
		a.ResponseType = ResponseMediaInfo
		return NewMediaInfo(a.MediaInfo)
	}

	return nil
}
