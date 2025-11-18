package structs

// DashboardResponse struct for response data dashboard
type DashboardResponse struct {
	CategoriesCount					int64			`json:"categories_count"`
	MapsCount						int64			`json:"maps_count"`
	ActiveMapsCount					int64			`json:"active_maps_count"`
	InactiveMapsCount				int64			`json:"inactive_maps_count"`
}