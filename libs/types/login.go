package types

type Login string

func (val Login) String() string {
	return string(val)
}

func (val Login) Validate() bool {
	return len(val) != 0
}
