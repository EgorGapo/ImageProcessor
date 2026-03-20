package usecases

type AuthService interface {
	Register(login string, password string) error
	Login(login string, password string) (string, error)
	Auth(token string) (string, error)
}
