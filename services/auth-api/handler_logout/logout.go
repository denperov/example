package handler_logout

import (
	"errors"

	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	DeleteSession(sessionID types.SessionID) bool
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
	SessionID types.SessionID
}

type Result struct {
}

const SessionNotFoundError = "session not found"

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	if !h.Storage.DeleteSession(p.SessionID) {
		return nil, errors.New(SessionNotFoundError)
	}

	return &Result{}, nil
}
