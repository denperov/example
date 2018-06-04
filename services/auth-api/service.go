package main

import (
	"github.com/denperov/owm-task/libs"
	"github.com/denperov/owm-task/services/auth-api/handler_check"
	"github.com/denperov/owm-task/services/auth-api/handler_create"
	"github.com/denperov/owm-task/services/auth-api/handler_delete"
	"github.com/denperov/owm-task/services/auth-api/handler_login"
	"github.com/denperov/owm-task/services/auth-api/handler_logout"
	"github.com/denperov/owm-task/services/auth-api/mongo"
	"github.com/labstack/echo"
)

func main() {
	storage := mongo.NewStorage("mongodb://auth-db/auth")
	defer storage.Close()

	apiServer := libs.APIServer{
		Echo: echo.New(),
	}
	apiServer.POST("/create", handler_create.New(storage))
	apiServer.POST("/delete", handler_delete.New(storage))
	apiServer.POST("/login", handler_login.New(storage))
	apiServer.POST("/logout", handler_logout.New(storage))
	apiServer.GET("/check", handler_check.New(storage))

	apiServer.Start(":8080")
}
