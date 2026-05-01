package response

import (
	"encoding/json"

	generalpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	igpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/ig/gen/igpb/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type WrapperResponse string

const (
	WRAPPER_RESPONSE_IG WrapperResponse = "ig"
)

type record = map[string]any

func isError(dataBin record) bool {
	return dataBin["error"] != nil
}

func getIgWrapperResponseRequest(data []byte, wrapper record, end_point ResponseType) Responser {
	var (
		gralResponse Responser
		errorMsg     = &generalpbv1.ResponseError{}
	)

	switch {
	// error
	case isError(wrapper):
		switch v := wrapper["error"].(type) {
		case string:
			errorMsg.SetCode(401)
			errorMsg.SetMessage(v)
		default:
			protojson.Unmarshal(data, errorMsg)
			gralResponse = NewResponse(errorMsg, ResponseError)
		}
	case end_point == SentMessageResponse:
		respMsg := &igpbv1.InstagramMessageResponse{}
		protojson.Unmarshal(data, respMsg)
		gralResponse = NewResponse(respMsg, SentMessageResponse)
	case end_point == InfoAccountBusinessResponse:
		infoAccount := &igpbv1.InstagramInfoAccountBusinessResponse{}
		protojson.Unmarshal(data, infoAccount)
		gralResponse = NewResponse(infoAccount, InfoAccountBusinessResponse)
	case end_point == InstagramFieldContainerResponse:
		containerResponse := &igpbv1.InstagramFieldContainerResponse{}
		protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, containerResponse)
		gralResponse = NewResponse(containerResponse, InstagramFieldContainerResponse)
	// general response
	case end_point == ResponseUnknow:
		fallthrough
	default:
		unknow := &generalpbv1.UnknownResponse{}
		dataValue := make(map[string]*structpb.Value)
		json.Unmarshal(data, &dataValue)
		unknow.SetData(dataValue)
		gralResponse = NewResponse(unknow, ResponseUnknow)
	}

	return gralResponse
}

func WrapperResponseRequest(dataBin []byte, wr WrapperResponse, endPointResponse ...ResponseType) Responser {
	var (
		end_point = ""
		errorMsg  = &generalpbv1.ResponseError{}
	)

	if len(endPointResponse) > 0 {
		end_point = endPointResponse[0]
	}
	wrapper := record{}

	err := json.Unmarshal(dataBin, &wrapper)
	if err != nil {
		errorMsg.SetCode(401)
		errorMsg.SetMessage(err.Error())
		return NewResponse(errorMsg, ResponseError)
	}

	//  generate response
	data, _ := json.Marshal(wrapper)

	switch wr {
	case WRAPPER_RESPONSE_IG:
		return getIgWrapperResponseRequest(data, wrapper, end_point)
	default:
		errorMsg.SetCode(401)
		errorMsg.SetMessage("wrapper response not found")
		return NewResponse(errorMsg, ResponseError)
	}
}
