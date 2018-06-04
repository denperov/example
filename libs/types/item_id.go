package types

type ItemID string

func NewItemID() ItemID {
	return ItemID(newID(16))
}

func (val ItemID) String() string {
	return string(val)
}

func (val ItemID) Validate() bool {
	return validateID(val, 16)
}
