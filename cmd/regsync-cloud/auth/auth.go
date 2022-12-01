package auth

type Authenticator interface {
	CheckRepo(repo string) (string, error)
	GetToken() string
}
