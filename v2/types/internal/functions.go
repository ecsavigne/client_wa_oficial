package internal

import "encoding/json"

var FirstNotEmpty = func(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}

	return ""
}

func ConvertStr(obj any) string {
	jsonValue, _ := json.Marshal(obj)

	return string(jsonValue)
}

func cleanDataEmptyMap(data map[string]any) map[string]any {
	for k, v := range data {
		if v == nil || v == "" || v == 0 || v == false {
			delete(data, k)
		} else if m, ok := v.(map[string]any); ok {
			data[k] = cleanDataEmptyMap(m)
		}
	}

	return data
}

func CleanDataEmpty(byteJson []byte) []byte {
	var data map[string]any
	err := json.Unmarshal(byteJson, &data)
	if err != nil {
		return byteJson
	}

	dataTemp := cleanDataEmptyMap(data)
	result, err := json.Marshal(dataTemp)
	if err != nil {
		return byteJson
	}

	return result
}
