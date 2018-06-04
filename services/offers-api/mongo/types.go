package mongo

import "github.com/denperov/owm-task/libs/types"

type TransferOfferCode struct {
	TransferOffer    *types.TransferOffer
	ConfirmationCode types.ConfirmationCode
}
