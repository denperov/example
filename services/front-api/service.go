package main

import (
	"fmt"
	"net/http"

	"github.com/denperov/owm-task/libs"
	"github.com/denperov/owm-task/libs/clients"
	"github.com/denperov/owm-task/libs/log"
	"github.com/denperov/owm-task/libs/types"
	"github.com/labstack/echo"
)

var (
	InvalidParamsError = echo.NewHTTPError(http.StatusBadRequest, "invalid params")
	UnauthorizedError  = echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
)

type Service struct {
	AuthAPIClient   *clients.AuthAPIClient
	ItemsAPIClient  *clients.ItemsAPIClient
	OffersAPIClient *clients.OffersAPIClient
}

func (s *Service) Start(address string) {
	e := echo.New()
	e.POST("/accounts/create", s.accountsCreate)
	e.POST("/accounts/login", s.accountsLogin)

	// API доступно только аутентифицированным пользователя
	apiWhitAuth := e.Group("", s.middleware())
	apiWhitAuth.POST("/accounts/delete", s.accountsDelete)
	apiWhitAuth.POST("/accounts/logout", s.accountsLogout)
	apiWhitAuth.GET("/items/list", s.itemsList)
	apiWhitAuth.POST("/items/create", s.itemsCreate)
	apiWhitAuth.POST("/items/delete", s.itemsDelete)
	apiWhitAuth.POST("/offers/create", s.offersCreate)
	apiWhitAuth.POST("/offers/confirm", s.offersConfirm)

	e.Logger.Fatal(e.Start(address))
}

func (s *Service) accountsCreate(c echo.Context) error {
	var params struct {
		Login    types.Login
		Password types.Password
	}
	if err := c.Bind(&params); err != nil {
		log.Warn(err)
		return InvalidParamsError
	}
	err := s.AuthAPIClient.Create(params.Login, params.Password)
	if err != nil {
		return s.apiError(err)
	}
	token, err := s.AuthAPIClient.Login(params.Login, params.Password)
	if err != nil {
		return s.apiError(err)
	}
	s.setCookieToken(token, c)
	return nil
}

func (s *Service) accountsLogin(c echo.Context) error {
	var params struct {
		Login    types.Login
		Password types.Password
	}
	if err := c.Bind(&params); err != nil {
		log.Warn(err)
		return InvalidParamsError
	}
	token, err := s.AuthAPIClient.Login(params.Login, params.Password)
	if err != nil {
		return s.apiError(err)
	}
	s.setCookieToken(token, c)
	return nil
}

func (s *Service) accountsDelete(c echo.Context) error {
	err := s.AuthAPIClient.Logout(SessionID(c))
	if err != nil {
		return s.apiError(err)
	}
	err = s.AuthAPIClient.Delete(UserID(c))
	if err != nil {
		return s.apiError(err)
	}
	s.resetCookieToken(c)
	return nil
}

func (s *Service) accountsLogout(c echo.Context) error {
	err := s.AuthAPIClient.Logout(SessionID(c))
	if err != nil {
		return s.apiError(err)
	}
	s.resetCookieToken(c)
	return nil
}

func (s *Service) itemsList(c echo.Context) error {
	items, err := s.ItemsAPIClient.List(UserID(c))
	if err != nil {
		return s.apiError(err)
	}
	return s.apiResult(items, c)
}

func (s *Service) itemsCreate(c echo.Context) error {
	var params struct {
		Content string
	}
	if err := c.Bind(&params); err != nil {
		log.Warn(err)
		return InvalidParamsError
	}
	itemID, err := s.ItemsAPIClient.Create(UserID(c), params.Content)
	if err != nil {
		return s.apiError(err)
	}
	return s.apiResult(struct {
		ItemID types.ItemID
	}{
		ItemID: itemID,
	}, c)
}

func (s *Service) itemsDelete(c echo.Context) error {
	var params struct {
		ItemID types.ItemID
	}
	if err := c.Bind(&params); err != nil {
		log.Warn(err)
		return InvalidParamsError
	}
	err := s.ItemsAPIClient.Delete(UserID(c), params.ItemID)
	if err != nil {
		return s.apiError(err)
	}
	return nil
}

func (s *Service) offersCreate(c echo.Context) error {
	var params struct {
		ItemID     types.ItemID
		NewOwnerID types.Login
	}
	if err := c.Bind(&params); err != nil {
		log.Warn(err)
		return InvalidParamsError
	}
	code, err := s.OffersAPIClient.TransferCreate(UserID(c), params.NewOwnerID, params.ItemID)
	if err != nil {
		return s.apiError(err)
	}
	return s.apiResult(struct {
		ConfirmationCode types.ConfirmationCode
	}{
		ConfirmationCode: code,
	}, c)
}

func (s *Service) offersConfirm(c echo.Context) error {
	var params struct {
		ConfirmationCode types.ConfirmationCode
	}
	if err := c.Bind(&params); err != nil {
		log.Warn(err)
		return InvalidParamsError
	}
	err := s.OffersAPIClient.TransferConfirm(UserID(c), params.ConfirmationCode)
	if err != nil {
		return s.apiError(err)
	}
	return nil
}

func (s *Service) apiOK() error {
	return nil
}

func (s *Service) apiResult(result interface{}, c echo.Context) error {
	return c.JSON(http.StatusOK, result)
}

func (s *Service) apiError(message interface{}) error {
	switch val := message.(type) {
	case string:
		return echo.NewHTTPError(http.StatusBadRequest, val)
	case fmt.Stringer:
		return echo.NewHTTPError(http.StatusBadRequest, val.String())
	case error:
		return echo.NewHTTPError(http.StatusBadRequest, val.Error())
	default:
		return echo.NewHTTPError(http.StatusBadRequest, val)
	}
}

func (s *Service) middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := s.getCookieToken(c)
			if !ok {
				return echo.ErrUnauthorized
			}
			session, err := s.AuthAPIClient.Check(token)
			if err != nil {
				log.Warn(err)
				return echo.ErrUnauthorized
			}
			return next(&ContextWithAuth{
				Context:   c,
				userID:    session.Login,
				sessionID: session.SessionID,
			})
		}
	}
}

func (s *Service) getCookieToken(c echo.Context) (types.Token, bool) {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return types.Token(""), false
	}
	token := types.Token(cookie.Value)
	if !token.Validate() {
		return types.Token(""), false
	}
	return token, true
}

func (s *Service) setCookieToken(token types.Token, c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:  "auth_token",
		Value: token.String(),
		Path:  "/",
	})
}

func (s *Service) resetCookieToken(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:  "auth_token",
		Value: "x",
		Path:  "/",
	})
}

type ContextWithAuth struct {
	echo.Context
	userID    types.Login
	sessionID types.SessionID
}

func (c *ContextWithAuth) UserID() types.Login {
	return c.userID
}

func (c *ContextWithAuth) SessionID() types.SessionID {
	return c.sessionID
}

func UserID(c echo.Context) types.Login {
	return c.(interface{ UserID() types.Login }).UserID()
}

func SessionID(c echo.Context) types.SessionID {
	return c.(interface{ SessionID() types.SessionID }).SessionID()
}

func main() {
	authService := &clients.AuthAPIClient{
		APIClient: libs.NewAPIClient("http://auth-api:8080"),
	}
	itemsService := &clients.ItemsAPIClient{
		APIClient: libs.NewAPIClient("http://items-api:8080"),
	}
	offersService := &clients.OffersAPIClient{
		APIClient: libs.NewAPIClient("http://offers-api:8080"),
	}

	service := Service{
		AuthAPIClient:   authService,
		ItemsAPIClient:  itemsService,
		OffersAPIClient: offersService,
	}
	service.Start(":8080")
}
