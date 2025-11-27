package structs

type MapResponse struct {
	Id					uint				`json:"id"`
	Image				string				`json:"image"`
	Name				string				`json:"name"`
	Slug				string				`json:"slug"`
	Description			string				`json:"description"`
	Address				string				`json:"address"`
	Latitude			string				`json:"latitude"`
	Longitude			string				`json:"longitude"`
	Geometry			string				`json:"geometry"`
	Status				string				`json:"status"`
	CategoryID			uint				`json:"category_id"`
	Category			struct {
		Id		uint		`json:"id"`
		Name	string		`json:"name"`
		Slug	string		`json:"slug"`
	} `json:"category"`
	CreatedAt			string				`json:"created_at"`
	UpdatedAt			string				`json:"updated_at"`
}

type PublicMap struct {
	Id					uint				`json:"id"`
	Name				string				`json:"name"`
	Slug				string				`json:"slug"`
	Description			string				`json:"description"`
	Address				string				`json:"address"`
	Latitude			string				`json:"latitude"`
	Longitude			string				`json:"longitude"`
	Geometry			string				`json:"geometry"` // Send string GeoJSON (like in DB)
	Status				string				`json:"status"` // "active" / "inactive"
	Image				string				`json:"image"` // file name
	CategoryID			uint				`json:"category_id"`
}

type MapCreateRequest struct {
	Name				string				`form:"name" json:"name" binding:"required,min=3"`
	Description			string				`form:"description" json:"description" binding:"required"`
	Address				string				`form:"address" json:"address" binding:"required"`
	Latitude			string				`form:"latitude" json:"latitude"`
	Longitude			string				`form:"longitude" json:"longitude"`
	Geometry			string				`form:"geometry" json:"geometry" binding:"omitempty"` // JSON string (GeoJSON / WKT which is serialized)
	Status				string				`form:"status" json:"status" binding:"omitempty,oneof=active inactive"` // default DB 'active' if empty
	CategoryID			uint				`form:"category_id" json:"category_id" binding:"required"`
}

type MapUpdateRequest struct {
	Name				string				`form:"name" json:"name" binding:"required,min=3"`
	Description			string				`form:"description" json:"description" bindng:"required"`
	Address				string				`form:"address" json:"address" binding:"required"`
	Latitude			string				`form:"latitude" json:"latitude"`
	Longitude			string				`form:"longitude" json:"longitude"`
	Geometry			string				`form:"geometry" json:"geometry" binding:"omitempty"`
	Status				string				`form:"status" json:"status" binding:"omitempty,oneof=active inactive"`
	CategoryID			uint				`form:"category_id" json:"category_id" binding:"required"`
}