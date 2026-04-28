package types

import (
	"fmt"
	"strings"

	"github.com/ecsavigne/client_wa_oficial/v2/types/internal"
)

func GetPermission() internal.PERMISSION_TYPE {
	return internal.GetPermissionType()
}

type QueryData map[string]any

func NewQueryData() QueryData {
	return make(map[string]any)
}

func (q QueryData) String() string {
	query := ""
	for k, v := range q {
		query += fmt.Sprintf("%s=%v&", k, v)
	}
	query = strings.TrimSuffix(query, "&")
	return query
}

func (q *QueryData) SetValue(key string, value any) {
	(*q)[key] = value
}

func (q QueryData) GetValue(key string) any {
	return q[key]
}
