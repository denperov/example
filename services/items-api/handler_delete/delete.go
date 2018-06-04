package handler_delete

import (
	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	DeleteUserItem(ownerID types.Login, itemID types.ItemID)
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

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	h.Storage.DeleteUserItem(p.OwnerID, p.ItemID)

	return &Result{}, nil
}
