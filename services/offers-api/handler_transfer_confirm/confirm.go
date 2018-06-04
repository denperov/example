package handler_transfer_confirm

import (
	"errors"

	"github.com/denperov/owm-task/libs/log"
	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	GetTransferOffer(recipientID types.Login, code types.ConfirmationCode) *types.TransferOffer
	DeleteTransferOffer(recipientID types.Login, code types.ConfirmationCode) bool
}

type ItemsAPI interface {
	Transfer(userID types.Login, itemID types.ItemID, newOwnerID types.Login) error
}

type handler struct {
	Storage  Storage
	ItemsAPI ItemsAPI
}

func New(storage Storage, itemsAPI ItemsAPI) *handler {
	return &handler{Storage: storage, ItemsAPI: itemsAPI}
}

func (h *handler) Params() interface{} {
	return &Params{}
}

type Params struct {
	RecipientID      types.Login
	ConfirmationCode types.ConfirmationCode
}

type Result struct {
}

const (
	OfferNotFoundError = "offer not found"
	ItemNotFoundError  = "item not found"
)

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	offer := h.Storage.GetTransferOffer(p.RecipientID, p.ConfirmationCode)
	if offer == nil {
		return nil, errors.New(OfferNotFoundError)
	}
	err := h.ItemsAPI.Transfer(offer.SenderID, offer.ItemID, offer.RecipientID)
	if err != nil {
		log.Error(err)
		return nil, errors.New(ItemNotFoundError)
	}
	h.Storage.DeleteTransferOffer(p.RecipientID, p.ConfirmationCode)
	return &Result{}, nil
}
