package auth

type Auther interface {
	CreateAccessToken() (string, error)
}

type Auth struct{}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) CreateAccessToken() (string, error) {
	// TODO: この関数内部のロジックは適宜変更する
	return "token", nil
}
