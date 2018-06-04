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
	err := s.items().EnsureIndex(mgo.Index{
		Key:        []string{"item.id"},
		Unique:     true,
		Background: true,
	})
	if err != nil {
		log.Panic(err)
	}

	err = s.items().EnsureIndex(mgo.Index{
		Key:        []string{"ownerid", "item.id"},
		Background: true,
	})
	if err != nil {
		log.Panic(err)
	}

}

func (s *Storage) items() *mgo.Collection {
	return s.Session.DB("").C("items")
}

func (s *Storage) CreateUserItem(ownerID types.Login, item *types.Item) {
	val := ItemOwnerID{
		OwnerID: ownerID,
		Item:    item,
	}
	err := s.items().Insert(&val)
	if err != nil {
		log.Panic(err)
	}
}

func (s *Storage) DeleteUserItem(ownerID types.Login, itemID types.ItemID) bool {
	err := s.items().Remove(bson.M{"ownerid": ownerID, "item.id": itemID})
	switch err {
	case nil:
		return true
	case mgo.ErrNotFound:
		return false
	default:
		panic(err)
	}
}

func (s *Storage) GetUserItems(ownerID types.Login) []*types.Item {
	var docs []*ItemOwnerID
	err := s.items().Find(bson.M{"ownerid": ownerID}).Select(bson.M{"item": 1}).All(&docs)
	if err != nil {
		log.Panic(err)
	}
	items := make([]*types.Item, 0)
	for _, doc := range docs {
		items = append(items, doc.Item)
	}
	return items
}

func (s *Storage) ChangeUserItemOwner(ownerID types.Login, itemID types.ItemID, newOwnerID types.Login) bool {
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{"ownerid": newOwnerID}},
	}
	_, err := s.items().Find(bson.M{"ownerid": ownerID, "item.id": itemID}).Apply(change, &ItemOwnerID{})
	switch err {
	case nil:
		return true
	case mgo.ErrNotFound:
		return false
	default:
		panic(err)
	}
}
