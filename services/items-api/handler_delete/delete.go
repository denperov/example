package handler_delete

import (
	"errors"

	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	DeleteUserItem(ownerID types.Login, itemID types.ItemID) bool
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
	OwnerID types.Login
	ItemID  types.ItemID
}

type Result struct {
}

const ItemNotFoundError = "item not found"

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	if !h.Storage.DeleteUserItem(p.OwnerID, p.ItemID) {
		return nil, errors.New(ItemNotFoundError)
	}

	return &Result{}, nil
}
