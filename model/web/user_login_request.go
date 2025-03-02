package web

type UserLoginRequest struct {
	Email     string `validate:"required,email" json:"email"`
	Password  string `validate:"required,min=6,max=100" json:"password"`
	UserAgent string `validate:"required" json:"user_agent"`
}
