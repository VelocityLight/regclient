package auth

import "fmt"

type Authenticator interface {
	CheckRepo(repo string) (string, error)
	GetToken() string
	GetLoginPassword() (string, string)
}

type commonCRAuth struct {
	Reg string
}

func NewCommonCRAuth(reg string) (Authenticator, error) {
	return commonCRAuth{Reg: reg}, nil
}

func (c commonCRAuth) CheckRepo(repo string) (string, error) {
	return fmt.Sprintf("%s/%s", c.Reg, repo), nil
}

func (c commonCRAuth) GetToken() string {
	return ""
}

func (c commonCRAuth) GetLoginPassword() (string, string) {
	return "", ""
}
