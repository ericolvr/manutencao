package dto

type UserRequest struct {
	Name     string `json:"name" validate:"required,min=5"`
	Mobile   string `json:"mobile" validate:"required,min=11,max=11"`
	Password string `json:"password" validate:"required,min=6"`
	Role     int64  `json:"role" validate:"required"`
	Status   bool   `json:"status"`
}

type UserResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Role   int64  `json:"role"`
	Status bool   `json:"status"`
}

type UserLogin struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Name  string `json:"name"`
	Token string `json:"token"`
	Role  int64  `json:"role"`
}
