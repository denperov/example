package handler_create

import (
	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	CreateUserItem(ownerID types.Login, item *types.Item)
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
	Content string
}

type Result struct {
	ItemID types.ItemID
}

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	item := types.Item{
		ID:      types.NewItemID(),
		Content: p.Content,
	}
	h.Storage.CreateUserItem(p.OwnerID, &item)

	return &Result{ItemID: item.ID}, nil
}
