package types

import (
	"github.com/denperov/owm-task/libs/log"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash string

func NewPasswordHash(password Password) PasswordHash {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password.String()), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}
	return PasswordHash(bytes)
}

func (hash PasswordHash) IsEmpty() bool {
	return len(hash) == 0
}

func (hash PasswordHash) Verify(password Password) bool {
	return !hash.IsEmpty() && bcrypt.CompareHashAndPassword([]byte(hash), []byte(password.String())) == nil
}
