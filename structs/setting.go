package structs

type SettingUpdateRequest struct {
	Title				string				`json:"title" binding:"required"`
	Description			string				`json:"description"`
	MapCenterLat		string				`json:"map_center_lat" binding:"required"`
	MapCenterLng		string				`json:"map_center_lng" binding:"required"`
	MapZoom				int					`json:"map_zoom" binding:"required,min=0,max=22"`
	VillageBoundary		string				`json:"village_boundary"`
}