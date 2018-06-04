package handler_delete

import (
	"errors"

	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	DeletePasswordHash(login types.Login) bool
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
	Login types.Login
}

type Result struct {
}

const LoginNotFoundError = "login not found"

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	ok := h.Storage.DeletePasswordHash(p.Login)
	if !ok {
		return nil, errors.New(LoginNotFoundError)
	}
	return &Result{}, nil
}
