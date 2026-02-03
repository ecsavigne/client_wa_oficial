package response

import (
	"encoding/json"
	"fmt"

	"github.com/ecsavigne/client_wa_oficial/v2/types"
)

type END_POINT string

const (
	// END_POINT_MEDIA_INFO    END_POINT = "media_info"
	END_POINT_BUSINESS_INFO END_POINT = "business_info"
	END_POINT_WABA_INFO     END_POINT = "waba_info"
	END_POINT_PHONE_INFO    END_POINT = "phone_info"
)

type record = map[string]any
type array = []any

func Val(r any) string {
	by, msg_e := json.MarshalIndent(r, "", "  ")
	if msg_e != nil {
		panic(fmt.Sprintf("Error occurred in MarshalIndent. Error is: %s", msg_e))
	}
	return string(by)
}

func isError(dataBin record) bool {
	return dataBin["error"] != nil
}

func isSuccess(wrapper record) bool {
	if wrapper["contacts"] != nil && wrapper["messaging_product"] != nil && wrapper["messages"] != nil {
		return true
	}
	return false
}

func isPhonesWA(wrapper record) bool {
	if wrapper["data"] != nil && wrapper["paging"] != nil && wrapper["data"].(array) != nil && wrapper["data"].(array)[0] != nil && wrapper["data"].(array)[0].(record)["quality_rating"] != nil && wrapper["data"].(array)[0].(record)["throughput"] != nil && wrapper["data"].(array)[0].(record)["id"] != nil {
		return true
	}
	return false
}

func isPhone(wrapper record) bool {
	if wrapper["id"] != nil && wrapper["verified_name"] != nil && wrapper["display_phone_number"] != nil && wrapper["code_verification_status"] != nil && wrapper["platform_type"] != nil && wrapper["quality_rating"] != nil {
		return true
	}
	return false
}

func isWaba(wrapper record) bool {
	if wrapper["id"] != nil && wrapper["name"] != nil /*&& wrapper["currency"] != nil*/ && wrapper["message_template_namespace"] != nil && wrapper["timezone_id"] != nil {
		return true
	}
	return false
}

func isTemplate(wrapper record) (bool, isOnly bool) {
	_, fCategory := wrapper["category"]
	_, fName := wrapper["name"]
	_, fLang := wrapper["language"]

	// only template
	if fCategory && fName && fLang {
		return true, true
	}

	data := wrapper["data"]
	if v, ok := data.(array); ok {
		if tpl, ok := v[0].(record); ok {
			_, fCategory = tpl["category"]
			_, fName = tpl["name"]
			_, fLang = tpl["language"]
			// array template
			if fCategory && fName && fLang {
				return true, false
			}
		}
	}

	return false, false
}

func isWabas(wrapper record) bool {
	if wrapper["data"] != nil && wrapper["paging"] != nil && wrapper["data"].(array) != nil &&
		wrapper["data"].(array)[0].(record)["id"] != nil && wrapper["data"].(array)[0].(record)["name"] != nil {
		return true
	}

	return false
}

func isBusiness(wrapper record) bool {
	if wrapper["name"] != nil && wrapper["two_factor_type"] != nil && wrapper["payment_account_id"] != nil && wrapper["is_hidden"] != nil ||
		wrapper["name"] != nil && wrapper["link"] != nil && wrapper["payment_account_id"] != nil && wrapper["is_hidden"] != nil {
		return true
	}

	return false
}

func isMediaInfo(wrapper record) bool {
	_, ok := wrapper["category"]
	if (wrapper["messaging_product"] != nil && wrapper["mime_type"] != nil && wrapper["id"] != nil && wrapper["url"] != nil && wrapper["sha256"] != nil) || (wrapper["id"] != nil && len(wrapper) <= 3 && !ok) {
		return true
	}
	return false
}

// JsonWrapperResponseRequest is a function that wraps a given json data into a ResponseWrapper.
// It unmarshals the data into a record, and then checks if the map contains an "error" key.
// If it does, it returns a types.Error object with the error message and code 401.
// Otherwise, it marshals the map back into json and unmarshals it into a types.GeneralResponse object.
// Finally, return interface que represent of type of response.
func JsonWrapperResponseRequest(dataBin []byte, endPoint ...END_POINT) Responser {
	end_point := END_POINT("")
	if len(endPoint) > 0 {
		end_point = endPoint[0]
	}
	wrapper := record{}
	gralResponse := NewGeneralResponse(&GeneralResponse{})

	err := json.Unmarshal(dataBin, &wrapper)
	if err != nil {
		return NewError(&Error{
			Code:    401,
			Message: err.Error(),
		})
	}

	//  generate response
	data, _ := json.Marshal(wrapper)
	isTpl, isOnly := isTemplate(wrapper)

	switch {
	// error
	case isError(wrapper):
		switch v := wrapper["error"].(type) {
		case string:
			return NewError(&Error{
				Message: v,
				Code:    401,
			})
		default:
			data, _ = json.Marshal(wrapper["error"])
			errorResponse := NewError(&Error{})
			json.Unmarshal(data, errorResponse)
			gralResponse.Error = errorResponse
		}

	// success
	case isSuccess(wrapper):
		successResponse := NewSuccess(&Success{})
		json.Unmarshal(data, successResponse)
		gralResponse.Success = successResponse

	// phonesWA
	case isPhonesWA(wrapper):
		phonesWA := NewPhonesWA(&PhonesWA{})
		json.Unmarshal(data, phonesWA)
		gralResponse.PhonesWA = phonesWA

	// Phone
	// case isPhone(wrapper):
	case end_point == END_POINT_PHONE_INFO:
		phone := NewPhone(&Phone{})
		json.Unmarshal(data, phone)
		gralResponse.Phone = phone

	// waba
	// case isWaba(wrapper):
	case end_point == END_POINT_WABA_INFO:
		wabaInfo := WabaInfo{}
		json.Unmarshal(data, &wabaInfo)
		waba := NewWABA(&Waba{
			Data: []WabaInfo{wabaInfo},
		})
		gralResponse.Waba = waba

	// template
	case isTpl:
		templates := NewTemplateResponse(&TemplateResponse{})
		if isOnly {
			mt := types.MockupTemplate{}
			json.Unmarshal(data, &mt)
			templates.Data = []types.MockupTemplate{mt}
		} else {
			json.Unmarshal(data, templates)
		}
		gralResponse.TemplateResponse = templates
	// wabas
	case isWabas(wrapper):
		waba := NewWABA(&Waba{})
		json.Unmarshal(data, waba)
		gralResponse.Waba = waba

	// business
	// case isBusiness(wrapper):
	case end_point == END_POINT_BUSINESS_INFO:
		business := NewBusiness(&Business{})
		json.Unmarshal(data, business)
		gralResponse.Business = business

	// media info
	case isMediaInfo(wrapper):
		mediaInfo := NewMediaInfo(&MediaInfo{})
		json.Unmarshal(data, mediaInfo)
		gralResponse.MediaInfo = mediaInfo

	// general response
	default:
		if gralResponse.OtherResponse == nil {
			gralResponse.OtherResponse = make(OtherResponse)
		}
		gralResponse.OtherResponse = OtherResponse(wrapper)
	}

	return gralResponse.GetResponseType()
}
