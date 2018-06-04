package libs

import (
	"net/http"

	"github.com/denperov/owm-task/libs/log"
	"github.com/gorilla/schema"
	"github.com/labstack/echo"
)

var queryDecoder = schema.NewDecoder()

type Handler interface {
	Params() (params interface{})
	Execute(params interface{}) (result interface{}, err error)
}

type APIServer struct {
	Echo *echo.Echo
}

func (s *APIServer) GET(path string, handler Handler) {
	s.Echo.GET(path, func(c echo.Context) error {
		params := handler.Params()
		err := queryDecoder.Decode(params, c.QueryParams())
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return s.executeHandler(handler, params, c)
	})
}

func (s *APIServer) POST(path string, handler Handler) {
	s.Echo.POST(path, func(c echo.Context) error {
		params := handler.Params()
		err := c.Bind(params)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return s.executeHandler(handler, params, c)
	})
}

func (s *APIServer) executeHandler(handler Handler, params interface{}, c echo.Context) error {
	result, err := handler.Execute(params)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if result != nil {
		return c.JSON(http.StatusOK, result)
	} else {
		return c.JSON(http.StatusOK, "")
	}
}

func (s *APIServer) Start(address string) {
	s.Echo.Logger.Fatal(s.Echo.Start(address))
}
