package mongo

import "github.com/denperov/owm-task/libs/types"

type ItemOwnerID struct {
	Item    *types.Item
	OwnerID types.Login
}
