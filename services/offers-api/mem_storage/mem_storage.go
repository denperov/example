package mem_storage

import (
	"sync"

	"github.com/denperov/owm-task/libs/types"
)

type Storage struct {
	offers sync.Map
	update sync.Mutex
}

type StorageTransferOffer struct {
	types.TransferOffer
	ConfirmationCode types.ConfirmationCode
}

func (s *Storage) CreateTransferOffer(offer *types.TransferOffer, code types.ConfirmationCode) {
	s.update.Lock()
	defer s.update.Unlock()

	if _, ok := s.offers.Load(code); ok {
		panic("confirmation code conflict")
	}
	s.offers.Store(code, offer)
}

func (s *Storage) GetTransferOffer(recipientID types.Login, code types.ConfirmationCode) *types.TransferOffer {
	if val, ok := s.offers.Load(code); ok && val.(*types.TransferOffer).RecipientID == recipientID {
		return val.(*types.TransferOffer)
	}
	return nil
}
