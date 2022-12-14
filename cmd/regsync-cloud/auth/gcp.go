package auth

import "fmt"

type gcrAuth struct {
	Reg     string
	Project string
}

func NewGcrAuth(reg string, proj string) (Authenticator, error) {
	return gcrAuth{
		Reg:     reg,
		Project: proj,
	}, nil
}

func (c gcrAuth) CheckRepo(repo string) (string, error) {
	return fmt.Sprintf("%s/%s/%s", c.Reg, c.Project, repo), nil
}

func (c gcrAuth) GetToken() string {
	return ""
}

func (c gcrAuth) GetLoginPassword() (string, string) {
	return "", ""
}

type garAuth struct {
	Reg string
}

func NewGarAuth(reg string) (Authenticator, error) {
	return garAuth{Reg: reg}, nil
}

func (c garAuth) CheckRepo(repo string) (string, error) {
	return fmt.Sprintf("%s/%s", c.Reg, repo), nil
}

func (c garAuth) GetToken() string {
	return ""
}

func (c garAuth) GetLoginPassword() (string, string) {
	return "", ""
}
