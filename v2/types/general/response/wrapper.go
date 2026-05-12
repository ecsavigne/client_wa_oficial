package response

import (
	"encoding/json"
	"strings"

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

func fixData(data []byte) []byte {
	fixedData := strings.ReplaceAll(string(data), "+0000", "Z")
	return []byte(fixedData)
}

func getIgWrapperResponseRequest(data []byte, wrapper record, end_point ResponseType) Responser {
	var (
		gralResponse Responser
		errorMsg     = &generalpbv1.ResponseError{}
		err          error
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
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, respMsg)
		respMsg.SetPayload(data)
		gralResponse = NewResponse(respMsg, SentMessageResponse)
	case end_point == InfoAccountBusinessResponse:
		infoAccount := &igpbv1.InstagramInfoAccountBusinessResponse{}
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, infoAccount)
		infoAccount.SetPayload(data)
		gralResponse = NewResponse(infoAccount, InfoAccountBusinessResponse)
	case end_point == InstagramFieldContainerResponse:
		containerResponse := &igpbv1.InstagramFieldContainerResponse{}
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, containerResponse)
		containerResponse.SetPayload(data)
		gralResponse = NewResponse(containerResponse, InstagramFieldContainerResponse)
	case end_point == ResponseSuccess:
		successResponse := &generalpbv1.SuccessResponse{}
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, successResponse)
		successResponse.SetPayload(data)
		gralResponse = NewResponse(successResponse, ResponseSuccess)
	case end_point == InstagramCommentResponse:
		commentResponse := &igpbv1.InstagramCommentResponseMessage{}
		data = fixData(data)
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, commentResponse)
		commentResponse.SetPayload(data)
		gralResponse = NewResponse(commentResponse, InstagramCommentResponse)
	// general response
	case end_point == InstagramMetricInsightResponse:
		metricInsightResponse := &igpbv1.InstagramMetricInsightResponse{}
		data = fixData(data)
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, metricInsightResponse)
		metricInsightResponse.SetPayload(data)
		gralResponse = NewResponse(metricInsightResponse, InstagramMetricInsightResponse)
	case end_point == InstagramMetricResponse:
		metricResponse := &igpbv1.InstagramMetricResponse{}
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, metricResponse)
		metricResponse.SetPayload(data)
		gralResponse = NewResponse(metricResponse, InstagramMetricResponse)
	case end_point == InstagramConversationMessageResponse:
		conversationMessageResponse := &igpbv1.InstagramConversationMessageResponse{}
		data = fixData(data)
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, conversationMessageResponse)
		conversationMessageResponse.SetPayload(data)
		gralResponse = NewResponse(conversationMessageResponse, InstagramConversationMessageResponse)
	case end_point == ConversationMessageResponse:
		conversationMessage := &igpbv1.ConversationMessage{}
		data = fixData(data)
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, conversationMessage)
		conversationMessage.SetPayload(data)
		gralResponse = NewResponse(conversationMessage, ConversationMessageResponse)
	case end_point == InstagramListConversationResponse:
		listConversationResponse := &igpbv1.ConversationMessages{}
		data = fixData(data)
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, listConversationResponse)
		listConversationResponse.SetPayload(data)
		gralResponse = NewResponse(listConversationResponse, InstagramListConversationResponse)
	case end_point == ResponseUnknow:
		fallthrough
	default:
		unknow := &generalpbv1.UnknownResponse{}
		dataMap := make(map[string]any)
		err = json.Unmarshal(data, &dataMap)
		dataValue, _ := structpb.NewValue(dataMap)
		unknow.SetData(dataValue)
		unknow.SetPayload(data)
		gralResponse = NewResponse(unknow, ResponseUnknow)
	}

	if err != nil {
		errorMsg.SetCode(401)
		errorMsg.SetMessage(err.Error())
		return NewResponse(errorMsg, ResponseError)
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
