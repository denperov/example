package handler_login

import (
	"errors"

	"github.com/denperov/owm-task/libs/types"
)

type Storage interface {
	GetPasswordHash(login types.Login) (types.PasswordHash, bool)
	CreateSession(session *types.Session, token types.Token)
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
	Token types.Token
}

const InvalidLoginError = "invalid login"
const InvalidPasswordError = "invalid password"

func (h *handler) Execute(params interface{}) (interface{}, error) {
	p := params.(*Params)

	passwordHash, ok := h.Storage.GetPasswordHash(p.Login)
	if !ok {
		return nil, errors.New(InvalidLoginError)
	}
	if !passwordHash.Verify(p.Password) {
		return nil, errors.New(InvalidPasswordError)
	}

	token := types.NewToken()
	session := types.Session{
		SessionID: types.NewSessionID(),
		Login:     p.Login,
	}
	h.Storage.CreateSession(&session, token)

	return &Result{Token: token}, nil
}
