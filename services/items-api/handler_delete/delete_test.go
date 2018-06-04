package handler_delete

import (
	"testing"

	"github.com/denperov/owm-task/libs/types"
	"github.com/denperov/owm-task/services/items-api/handler_delete/mocks"
	"github.com/stretchr/testify/assert"
)

var GoodOwnerID = types.Login("good_login")
var GoodItemID = types.NewItemID()
var NotFoundItemID = types.NewItemID()

func TestOK(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("DeleteUserItem", GoodOwnerID, GoodItemID).Return(true)

	handler := New(storage)

	params := &Params{OwnerID: GoodOwnerID, ItemID: GoodItemID}
	result, err := handler.Execute(params)
	if assert.NoError(t, err) {
		assert.IsType(t, result, &Result{})
	}
}

func TestNotFound(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("DeleteUserItem", GoodOwnerID, NotFoundItemID).Return(false)

	handler := New(storage)

	params := &Params{OwnerID: GoodOwnerID, ItemID: NotFoundItemID}
	_, err := handler.Execute(params)
	assert.EqualError(t, err, ItemNotFoundError)
}
