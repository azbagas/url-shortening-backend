package domain

type RefreshToken struct {
	Id           int
	RefreshToken string
	UserId       int
	UserAgent    string
}
