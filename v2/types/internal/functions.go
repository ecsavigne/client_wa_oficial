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
