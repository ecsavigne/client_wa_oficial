package types

import "github.com/ecsavigne/client_wa_oficial/v2/types/internal"

func GetPermission() internal.PERMISSION_TYPE {
	return internal.GetPermissionType()
}
