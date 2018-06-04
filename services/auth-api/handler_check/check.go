package handler_check

import (
	"errors"

	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	CheckToken(token types.Token) *types.Session
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
	Token types.Token
}

type Result struct {
	Session *types.Session
}

const TokenNotFoundError = "token not found"

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	session := h.Storage.CheckToken(p.Token)
	if session == nil {
		return nil, errors.New(TokenNotFoundError)
	}

	return &Result{Session: session}, nil
}
