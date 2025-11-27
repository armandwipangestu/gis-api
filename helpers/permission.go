package helpers

import (
	"armandwipangestu/gis-api/models"
)

func GetPermissionMap(roles []models.Role) map[string]bool {
	permissionMap := make(map[string]bool)

	for _, role := range roles {
		for _, perm := range role.Permissions {
			permissionMap[perm.Name] = true
		}
	}

	return permissionMap
}