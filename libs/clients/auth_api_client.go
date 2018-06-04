package clients

import (
	"github.com/denperov/owm-task/libs"
	"github.com/denperov/owm-task/libs/types"
	"github.com/denperov/owm-task/services/auth-api/handler_check"
	"github.com/denperov/owm-task/services/auth-api/handler_create"
	"github.com/denperov/owm-task/services/auth-api/handler_delete"
	"github.com/denperov/owm-task/services/auth-api/handler_login"
	"github.com/denperov/owm-task/services/auth-api/handler_logout"
)

type AuthAPIClient struct {
	APIClient *libs.APIClient
}

func (c *AuthAPIClient) Check(token types.Token) (*types.Session, error) {
	params := handler_check.Params{Token: token}
	var result handler_check.Result

	err := c.APIClient.Get("/check", &params, &result)
	return result.Session, err
}

func (c *AuthAPIClient) Create(login types.Login, password types.Password) error {
	params := handler_create.Params{Login: login, Password: password}
	var result handler_create.Result

	return c.APIClient.Post("/create", &params, &result)
}

func (c *AuthAPIClient) Delete(login types.Login) error {
	params := handler_delete.Params{Login: login}
	var result handler_delete.Result

	return c.APIClient.Post("/delete", &params, &result)
}

func (c *AuthAPIClient) Login(login types.Login, password types.Password) (types.Token, error) {
	params := handler_login.Params{Login: login, Password: password}
	var result handler_login.Result

	err := c.APIClient.Post("/login", &params, &result)
	return result.Token, err
}

func (c *AuthAPIClient) Logout(sessionID types.SessionID) error {
	params := handler_logout.Params{SessionID: sessionID}
	var result handler_logout.Result

	return c.APIClient.Post("/logout", &params, &result)
}
