package structs

// This struct is used for show data user as API response
type UserResponse struct {
	Id				uint				`json:"id"`
	Name			string				`json:"name"`
	Username		string				`json:"username"`
	Email			string				`json:"email"`
	Permissions		map[string]bool		`json:"permissions,omitempty"` // Show permissions that owned of user
	Roles			[]RoleResponse		`json:"roles,omitempty"`
	Token			*string				`json:"token,omitempty"`
	CreatedAt		string				`json:"created_at,omitempty"`
	UpdatedAt		string				`json:"updated_at,omitempty"`
}

type UserSimpleResponse struct {
	Id				uint				`json:"id"`
	Name			string				`json:"name"`
}

// This struct is used for receive data and processing when creating user
type UserCreateRequest struct {
	Name			string				`json:"name" binding:"required"`
	Username		string				`json:"username" binding:"required" gorm:"unique;not null"`
	Email			string				`json:"email" binding:"required" gorm:"unique;not null"`
	Password		string				`json:"password" binding:"required"`
	RoleIDs			[]uint				`json:"role_ids" binding:"required"` // IDs of roles that assigned to the user
}

// This struct is used for receive data and processing when update user
type UserUpdateRequest struct {
	Name			string				`json:"name" binding:"required"`
	Username		string				`json:"username" binding:"required" gorm:"unique;not null"`
	Email			string				`json:"email" binding:"required" gorm:"unique;not null"`
	Password		string				`json:"password,omitempty"`
	RoleIDs			[]uint				`json:"role_ids"` // IDs of role that assigned to the user
}

// This struct is used when user processing to login
type UserLoginRequest struct {
	Username		string				`json:"username" binding:"required"`
	Password		string				`json:"password" binding:"required"`
}