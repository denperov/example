package mongo

import (
	"github.com/denperov/owm-task/libs/log"
	"github.com/denperov/owm-task/libs/types"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func NewStorage(address string) *Storage {
	s, err := mgo.Dial(address)
	if err != nil {
		log.Panic(err)
	}
	st := &Storage{
		Session: s,
	}
	st.itemsInit()
	return st
}

type Storage struct {
	Session *mgo.Session
}

func (s *Storage) Close() {
	s.Session.Close()
}

func (s *Storage) itemsInit() {
	err := s.offers().EnsureIndex(mgo.Index{
		Key:        []string{"confirmationcode"},
		Unique:     true,
		Background: true,
	})
	if err != nil {
		log.Panic(err)
	}

	err = s.offers().EnsureIndex(mgo.Index{
		Key:        []string{"confirmationcode", "transferoffer.recipientid"},
		Background: true,
	})
	if err != nil {
		log.Panic(err)
	}

}

func (s *Storage) offers() *mgo.Collection {
	return s.Session.DB("").C("offers")
}

func (s *Storage) CreateTransferOffer(offer *types.TransferOffer, code types.ConfirmationCode) {
	val := TransferOfferCode{
		TransferOffer:    offer,
		ConfirmationCode: code,
	}
	err := s.offers().Insert(&val)
	if err != nil {
		log.Panic(err)
	}
}

func (s *Storage) GetTransferOffer(recipientID types.Login, code types.ConfirmationCode) *types.TransferOffer {
	var result TransferOfferCode
	err := s.offers().Find(bson.M{"confirmationcode": code, "transferoffer.recipientid": recipientID}).One(&result)
	switch err {
	case nil:
		return result.TransferOffer
	case mgo.ErrNotFound:
		return nil
	default:
		panic(err)
	}
}

func (s *Storage) DeleteTransferOffer(recipientID types.Login, code types.ConfirmationCode) bool {
	err := s.offers().Remove(bson.M{"confirmationcode": code, "transferoffer.recipientid": recipientID})
	switch err {
	case nil:
		return true
	case mgo.ErrNotFound:
		return false
	default:
		panic(err)
	}
}
