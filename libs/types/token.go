package types

type Token string

func NewToken() Token {
	return Token(newID(24))
}

func (val Token) String() string {
	return string(val)
}

func (val Token) Validate() bool {
	return validateID(val, 24)
}
