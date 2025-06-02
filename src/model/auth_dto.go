package model

type SignupRequest struct {
	Email    string   `json:"email" binding:"required" example:"user@example.com"`
	Password string   `json:"password" binding:"required" example:"strongpassword"`
	Name     string   `json:"name" binding:"required" example:"John"`
	Surname  string   `json:"surname" binding:"required" example:"Doe"`
	Phone    string   `json:"phone" binding:"required" example:"+123456789"`
	Role     UserRole `json:"role" binding:"required" example:"editor"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"strongpassword"`
}

type UpdateUserRequest struct {
	Email   string   `json:"email" example:"user@example.com"`
	Name    string   `json:"name" example:"John"`
	Surname string   `json:"surname" example:"Doe"`
	Phone   string   `json:"phone" example:"+123456789"`
	Role    UserRole `json:"role" example:"editor"`
}

// Paginated list response DTOs


