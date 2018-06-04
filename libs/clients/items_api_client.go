package clients

import (
	"github.com/denperov/owm-task/libs"
	"github.com/denperov/owm-task/libs/types"
	"github.com/denperov/owm-task/services/items-api/handler_create"
	"github.com/denperov/owm-task/services/items-api/handler_delete"
	"github.com/denperov/owm-task/services/items-api/handler_list"
	"github.com/denperov/owm-task/services/items-api/handler_transfer"
)

type ItemsAPIClient struct {
	APIClient *libs.APIClient
}

func (c *ItemsAPIClient) Create(ownerID types.Login, content string) (types.ItemID, error) {
	params := handler_create.Params{OwnerID: ownerID, Content: content}
	var result handler_create.Result

	err := c.APIClient.Post("/create", &params, &result)
	return result.ItemID, err
}

func (c *ItemsAPIClient) Delete(ownerID types.Login, itemID types.ItemID) error {
	params := handler_delete.Params{OwnerID: ownerID, ItemID: itemID}
	var result handler_delete.Result

	return c.APIClient.Post("/delete", &params, &result)
}

func (c *ItemsAPIClient) List(ownerID types.Login) ([]*types.Item, error) {
	params := handler_list.Params{OwnerID: ownerID}
	var result handler_list.Result

	err := c.APIClient.Get("/list", &params, &result)
	return result.Items, err
}

func (c *ItemsAPIClient) Transfer(ownerID types.Login, itemID types.ItemID, newOwnerID types.Login) error {
	params := handler_transfer.Params{OwnerID: ownerID, ItemID: itemID, NewOwnerID: newOwnerID}
	var result handler_transfer.Result

	return c.APIClient.Post("/transfer", &params, &result)
}
