package web

type UserRegisterRequest struct {
	Name                 string `validate:"required,min=1,max=100" json:"name"`
	Email                string `validate:"required,email" json:"email"`
	Password             string `validate:"required,min=6,max=100" json:"password"`
	PasswordConfirmation string `validate:"required,min=6,max=100,eqfield=Password" json:"passwordConfirmation"`
}
