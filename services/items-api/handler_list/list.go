package handler_list

import (
	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	GetUserItems(ownerID types.Login) []*types.Item
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
}

type Result struct {
	Items []*types.Item
}

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	return &Result{
		Items: h.Storage.GetUserItems(p.OwnerID),
	}, nil
}
