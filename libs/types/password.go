package types

type Password string

func (val Password) String() string {
	return string(val)
}

func (val Password) Validate() bool {
	return len(val) != 0
}
