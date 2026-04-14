package response

import (
	"encoding/json"

	generalpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/general/gen/generalpb/v1"
	"github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	igpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/ig/gen/igpb/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type record = map[string]any

func isError(dataBin record) bool {
	return dataBin["error"] != nil
}

func WrapperResponseRequest(dataBin []byte, endPointResponse ...response.ResponseType) response.Responser {
	var (
		gralResponse response.Responser
		end_point    = ""
		errorMsg     = &generalpbv1.ResponseError{}
	)

	if len(endPointResponse) > 0 {
		end_point = endPointResponse[0]
	}
	wrapper := record{}

	err := json.Unmarshal(dataBin, &wrapper)
	if err != nil {
		errorMsg.SetCode(401)
		errorMsg.SetMessage(err.Error())
		gralResponse = response.NewResponse(errorMsg, response.ResponseError)
		return gralResponse
	}

	//  generate response
	data, _ := json.Marshal(wrapper)

	switch {
	// error
	case isError(wrapper):
		switch v := wrapper["error"].(type) {
		case string:
			errorMsg.SetCode(401)
			errorMsg.SetMessage(v)
		default:
			protojson.Unmarshal(data, errorMsg)
			gralResponse = response.NewResponse(errorMsg, response.ResponseError)
		}

	case end_point == response.ResponseUnknow:
		unknow := &generalpbv1.UnknownResponse{}
		protojson.Unmarshal(data, unknow)
		gralResponse = response.NewResponse(unknow, response.ResponseUnknow)

	case end_point == ResponseSentMessage:
		respMsg := &igpbv1.InstagramMessageResponse{}
		protojson.Unmarshal(data, respMsg)
		gralResponse = response.NewResponse(respMsg, ResponseSentMessage)

	// general response
	default:
		unknow := &generalpbv1.UnknownResponse{}
		protojson.Unmarshal(data, unknow)
		gralResponse = response.NewResponse(unknow, response.ResponseUnknow)
	}

	return gralResponse
}
