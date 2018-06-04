package mongo

import "github.com/denperov/owm-task/libs/types"

type LoginPasswordHash struct {
	Login        types.Login
	PasswordHash types.PasswordHash
}

type SessionToken struct {
	*types.Session
	Token types.Token
}
