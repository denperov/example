package types

type TransferOffer struct {
	SenderID    Login
	RecipientID Login
	ItemID      ItemID
}
