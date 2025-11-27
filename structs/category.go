package structs

type CategoryResponse struct {
	Id				uint			`json:"id"`
	Image			string			`json:"image"`
	Name			string			`json:"name"`
	Slug			string			`json:"slug"`
	Color			string			`json:"color"`
	Description		string			`json:"description"`
	Maps			[]PublicMap		`json:"maps"`
	CreatedAt		string			`json:"created_at"`
	UpdatedAt		string			`json:"updated_at"`
}

type CategoryCreateRequest struct {
	Name			string			`form:"name" bindng:"required"`
	Color			string			`form:"color" binding:"required"`
	Description		string			`form:"description" binding:"required"`
}

type CategoryUpdateRequest struct {
	Name			string			`form:"name" binding:"required"`
	Color			string			`form:"color" binding:"required"`
	Description		string			`form:"description" binding:"required"`
}