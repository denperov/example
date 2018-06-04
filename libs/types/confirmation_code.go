package types

type ConfirmationCode string

func NewConfirmationCode() ConfirmationCode {
	return ConfirmationCode(newID(16))
}

func (val ConfirmationCode) String() string {
	return string(val)
}

func (val ConfirmationCode) Validate() bool {
	return validateID(val, 16)
}
