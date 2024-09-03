// Package auth is used for HTTP authentication
package auth

import "fmt"

type acrAuth struct {
	Reg     string
}

func NewAcrAuth(reg string) (Authenticator, error) {
	return acrAuth{
		Reg:     reg,
	}, nil
}

func (c acrAuth) CheckRepo(repo string) (string, error) {
	return fmt.Sprintf("%s/%s", c.Reg, repo), nil
}

func (c acrAuth) GetToken() string {
	return ""
}

func (c acrAuth) GetLoginPassword() (string, string) {
	return "", ""
}