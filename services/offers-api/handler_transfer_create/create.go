package handler_transfer_create

import (
	"errors"

	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	CreateTransferOffer(offer *types.TransferOffer, code types.ConfirmationCode)
}

type handler struct {
	Storage Storage
}

func New(storage Storage) *handler {
	return &handler{Storage: storage}
}

func (h *handler) Params() interface{} {
	return &Params{}
}

type Params struct {
	SenderID    types.Login
	RecipientID types.Login
	ItemID      types.ItemID
}

type Result struct {
	ConfirmationCode types.ConfirmationCode
}

const InvalidRecipientError = "invalid recipient"

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	//TODO: Уточнить возможность создавать предложения для объектов, которыми пользователь не обладает

	if p.SenderID == p.RecipientID {
		return nil, errors.New(InvalidRecipientError)
	}

	code := types.NewConfirmationCode()
	offer := types.TransferOffer{
		SenderID:    p.SenderID,
		RecipientID: p.RecipientID,
		ItemID:      p.ItemID,
	}
	h.Storage.CreateTransferOffer(&offer, code)

	return &Result{ConfirmationCode: code}, nil
}
