package structs

type PermissionCreateRequest struct {
	Name	string	`json:"name" binding:"required"`
}

type PermissionUpdateRequest struct {
	Name	string 	`json:"name" binding:"required"`
}

type PermissionResponse struct {
	Id			uint	`json:"id"`
	Name		string	`json:"name"`
	CreatedAt 	string	`json:"created_at"`
	UpdatedAt	string	`json:"updated_at"`
}