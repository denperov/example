package types

import (
	"crypto/rand"
	"encoding/base32"

	"github.com/denperov/owm-task/libs/log"
)

func newID(n int) string {
	buff := make([]byte, base32.StdEncoding.DecodedLen(n))
	if _, err := rand.Read(buff); err != nil {
		log.Panic(err)
	}
	return base32.StdEncoding.EncodeToString(buff)
}

func validateID(val interface{}, n int) bool {
	str := val.(interface{ String() string }).String()
	if len(str) != n {
		return false
	}
	if _, err := base32.StdEncoding.DecodeString(str); err != nil {
		return false
	}
	return true
}
