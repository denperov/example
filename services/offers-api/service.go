package main

import (
	"github.com/denperov/owm-task/libs"
	"github.com/denperov/owm-task/libs/clients"
	"github.com/denperov/owm-task/services/offers-api/handler_transfer_confirm"
	"github.com/denperov/owm-task/services/offers-api/handler_transfer_create"
	"github.com/denperov/owm-task/services/offers-api/mongo"
	"github.com/labstack/echo"
)

func main() {
	storage := mongo.NewStorage("mongodb://offers-db/offers")
	defer storage.Close()

	itemsService := clients.ItemsAPIClient{
		APIClient: libs.NewAPIClient("http://items-api:8080"),
	}

	apiServer := libs.APIServer{
		Echo: echo.New(),
	}
	apiServer.POST("/transfer/create", handler_transfer_create.New(storage))
	apiServer.POST("/transfer/confirm", handler_transfer_confirm.New(storage, &itemsService))

	apiServer.Start(":8080")
}
