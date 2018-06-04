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
	st.passwordsInit()
	st.sessionsInit()
	return st
}

type Storage struct {
	Session *mgo.Session
}

func (s *Storage) Close() {
	s.Session.Close()
}

func (s *Storage) passwordsInit() {
	err := s.passwords().EnsureIndex(mgo.Index{
		Key:        []string{"login"},
		Unique:     true,
		Background: true,
	})
	if err != nil {
		log.Panic(err)
	}
}

func (s *Storage) sessionsInit() {
	err := s.sessions().EnsureIndex(mgo.Index{
		Key:        []string{"token"},
		Unique:     true,
		Background: true,
	})
	if err != nil {
		log.Panic(err)
	}
	err = s.sessions().EnsureIndex(mgo.Index{
		Key:        []string{"session.sessionid"},
		Unique:     true,
		Background: true,
	})
	if err != nil {
		log.Panic(err)
	}
}

func (s *Storage) passwords() *mgo.Collection {
	return s.Session.DB("").C("passwords")
}

func (s *Storage) sessions() *mgo.Collection {
	return s.Session.DB("").C("sessions")
}

func (s *Storage) CreatePasswordHash(login types.Login, hash types.PasswordHash) bool {
	val := LoginPasswordHash{
		Login:        login,
		PasswordHash: hash,
	}
	err := s.passwords().Insert(&val)
	if err != nil {
		if mgo.IsDup(err) {
			return false
		}
		log.Panic(err)
	}
	return true
}

func (s *Storage) GetPasswordHash(login types.Login) (types.PasswordHash, bool) {
	var result LoginPasswordHash
	err := s.passwords().Find(bson.M{"login": login}).One(&result)
	switch err {
	case nil:
		return result.PasswordHash, true
	case mgo.ErrNotFound:
		return types.PasswordHash(""), false
	default:
		panic(err)
	}
}

func (s *Storage) DeletePasswordHash(login types.Login) bool {
	err := s.passwords().Remove(bson.M{"login": login})
	switch err {
	case nil:
		return true
	case mgo.ErrNotFound:
		return false
	default:
		panic(err)
	}
}

func (s *Storage) CreateSession(session *types.Session, token types.Token) {
	val := SessionToken{
		Session: session,
		Token:   token,
	}
	err := s.sessions().Insert(&val)
	if err != nil {
		log.Panic(err)
	}
}

func (s *Storage) DeleteSession(sessionID types.SessionID) bool {
	err := s.sessions().Remove(bson.M{"session.sessionid": sessionID})
	switch err {
	case nil:
		return true
	case mgo.ErrNotFound:
		return false
	default:
		panic(err)
	}
}

func (s *Storage) CheckToken(token types.Token) *types.Session {
	var result SessionToken
	err := s.sessions().Find(bson.M{"token": token}).One(&result)
	switch err {
	case nil:
		return result.Session
	case mgo.ErrNotFound:
		return nil
	default:
		panic(err)
	}
}
