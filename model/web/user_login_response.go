package web

type UserLoginResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"accessToken"`
}

type UserLoginResponseWithRefreshToken struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}
