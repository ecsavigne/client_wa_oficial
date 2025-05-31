package response

import (
	"fmt"
	"io"

	"github.com/ecsavigne/client_wa_oficial/v2/types"
)

type GeneralResponse struct {
	KernelResponser
	ResponseType string `json:"response_type,omitempty"`
	*PhonesWA    `json:",omitempty"`
	*Phone       `json:",omitempty"`
	*Waba        `json:",omitempty"`
	*Error       `json:",omitempty"`
	*Success     `json:",omitempty"`
	*MediaInfo   `json:",omitempty"`
	*Business    `json:",omitempty"`
}

func NewGeneralResponse(config ResponserRequest) *GeneralResponse {
	if v, ok := config.(*GeneralResponse); ok {
		v.ResponseType = ResponseGeneralResponse
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not RespenserRequest")
}

func (a *GeneralResponse) GetResponseType() ResponserRequest {
	switch {
	case a.PhonesWA != nil:
		a.ResponseType = ResponsePhonesWA
		return NewPhonesWA(a.PhonesWA)
	case a.Phone != nil:
		a.ResponseType = ResponsePhone
		return NewPhone(a.Phone)
	case a.Waba != nil:
		a.ResponseType = ResponseWABA
		return NewWABA(a.Waba)
	case a.Error != nil:
		a.ResponseType = ResponseError
		return NewError(a.Error)
	case a.Success != nil:
		a.ResponseType = ResponseSuccess
		return NewSuccess(a.Success)
	case a.MediaInfo != nil:
		a.ResponseType = ResponseMediaInfo
		return NewMediaInfo(a.MediaInfo)
	case a.Business != nil:
		a.ResponseType = ResponseBusiness
		return NewBusiness(a.Business)
	default:
		return a
	}
}

func GetResponseRequest(bodyResponse io.ReadCloser, funcName, who string) ResponserRequest {
	b, err := io.ReadAll(bodyResponse)
	defer bodyResponse.Close()
	if err != nil {
		return NewError(&Error{
			Type:    ResponseError,
			Code:    types.CodeErrorUnrecognized,
			Message: fmt.Sprintf("Error reading in %s the data of request of %s. error is: %v", funcName, who, err),
		})
	}

	return JsonWrapperResponseRequest(b)
}
