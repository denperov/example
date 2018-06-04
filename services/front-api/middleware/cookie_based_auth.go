package middleware

import (
	"net/http"
	"time"

	"github.com/denperov/owm-task/libs/clients"
	"github.com/denperov/owm-task/libs/log"
	"github.com/denperov/owm-task/libs/types"
	"github.com/labstack/echo"
)

type CookieBasedAuth struct {
	AuthCookieName string
	AuthAPIClient  *clients.AuthAPIClient
}

func (m *CookieBasedAuth) GetCookieToken(c echo.Context) (types.Token, bool) {
	cookie, err := c.Cookie(m.AuthCookieName)
	if err != nil {
		return types.Token(""), false
	}
	token := types.Token(cookie.Value)
	if !token.Validate() {
		return types.Token(""), false
	}
	return token, true
}

func (m *CookieBasedAuth) SetCookieToken(token types.Token, c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:    m.AuthCookieName,
		Value:   token.String(),
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})
}

func (m *CookieBasedAuth) ResetCookieToken(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:  m.AuthCookieName,
		Value: "x",
		Path:  "/",
	})
}

func (m *CookieBasedAuth) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := m.GetCookieToken(c)
			if !ok {
				return echo.ErrUnauthorized
			}
			session, err := m.AuthAPIClient.Check(token)
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
