package main

import (
	"github.com/denperov/owm-task/libs"
	"github.com/denperov/owm-task/services/items-api/handler_create"
	"github.com/denperov/owm-task/services/items-api/handler_delete"
	"github.com/denperov/owm-task/services/items-api/handler_list"
	"github.com/denperov/owm-task/services/items-api/handler_transfer"
	"github.com/denperov/owm-task/services/items-api/mongo"
	"github.com/labstack/echo"
)

func main() {
	storage := mongo.NewStorage("mongodb://items-db/items")
	defer storage.Close()

	apiServer := libs.APIServer{
		Echo: echo.New(),
	}
	apiServer.GET("/list", handler_list.New(storage))
	apiServer.POST("/create", handler_create.New(storage))
	apiServer.POST("/delete", handler_delete.New(storage))
	apiServer.POST("/transfer", handler_transfer.New(storage))

	apiServer.Start(":8080")
}
