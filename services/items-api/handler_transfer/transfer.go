package handler_transfer

import (
	"errors"

	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	ChangeUserItemOwner(ownerID types.Login, itemID types.ItemID, newOwnerID types.Login) bool
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
	OwnerID    types.Login
	ItemID     types.ItemID
	NewOwnerID types.Login
}

type Result struct {
}

const ItemNotFoundError = "item not found"

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	if !h.Storage.ChangeUserItemOwner(p.OwnerID, p.ItemID, p.NewOwnerID) {
		return nil, errors.New(ItemNotFoundError)
	}
	return &Result{}, nil
}
