package types

type SessionID string

func NewSessionID() SessionID {
	return SessionID(newID(16))
}

func (val SessionID) String() string {
	return string(val)
}

func (val SessionID) Validate() bool {
	return validateID(val, 16)
}
