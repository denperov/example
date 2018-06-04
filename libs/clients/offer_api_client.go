package clients

import (
	"github.com/denperov/owm-task/libs"
	"github.com/denperov/owm-task/libs/types"
	"github.com/denperov/owm-task/services/offers-api/handler_transfer_confirm"
	"github.com/denperov/owm-task/services/offers-api/handler_transfer_create"
)

type OffersAPIClient struct {
	APIClient *libs.APIClient
}

func (c *OffersAPIClient) TransferConfirm(recipientID types.Login, code types.ConfirmationCode) error {
	params := handler_transfer_confirm.Params{RecipientID: recipientID, ConfirmationCode: code}
	var result handler_transfer_confirm.Result

	return c.APIClient.Post("/transfer/confirm", &params, &result)
}

func (c *OffersAPIClient) TransferCreate(senderID types.Login, recipientID types.Login, itemID types.ItemID) (types.ConfirmationCode, error) {
	params := handler_transfer_create.Params{SenderID: senderID, RecipientID: recipientID, ItemID: itemID}
	var result handler_transfer_create.Result

	err := c.APIClient.Post("/transfer/create", &params, &result)
	return result.ConfirmationCode, err
}
