package handler_create

import (
	"errors"

	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	CreatePasswordHash(login types.Login, hash types.PasswordHash) bool
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
	Login    types.Login
	Password types.Password
}

type Result struct {
}

const LoginAlreadyExistsError = "login already exists"

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	if !h.Storage.CreatePasswordHash(p.Login, types.NewPasswordHash(p.Password)) {
		return nil, errors.New(LoginAlreadyExistsError)
	}

	return &Result{}, nil
}
