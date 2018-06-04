package handler_transfer

import (
	"testing"

	"github.com/denperov/owm-task/libs/types"
	"github.com/denperov/owm-task/services/items-api/handler_transfer/mocks"
	"github.com/stretchr/testify/assert"
)

var GoodOwnerID = types.Login("good_login")
var GoodNewOwnerID = types.Login("good_new_login")
var NotFoundOwnerID = types.Login("not_found_login")
var NotFoundNewOwnerID = types.Login("not_found_new_login")
var GoodItemID = types.NewItemID()
var NotFoundItemID = types.NewItemID()

func TestOK(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("ChangeUserItemOwner", GoodOwnerID, GoodItemID, GoodNewOwnerID).Return(true)

	handler := New(storage)

	params := &Params{OwnerID: GoodOwnerID, ItemID: GoodItemID, NewOwnerID: GoodNewOwnerID}
	result, err := handler.Execute(params)
	if assert.NoError(t, err) {
		assert.IsType(t, result, &Result{})
	}
}

func TestOwnerNotFound(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("ChangeUserItemOwner", NotFoundOwnerID, GoodItemID, GoodNewOwnerID).Return(true)

	handler := New(storage)

	params := &Params{OwnerID: NotFoundOwnerID, ItemID: GoodItemID, NewOwnerID: GoodNewOwnerID}
	result, err := handler.Execute(params)
	if assert.NoError(t, err) {
		assert.IsType(t, result, &Result{})
	}
}

func TestNewOwnerNotFound(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("ChangeUserItemOwner", GoodOwnerID, GoodItemID, NotFoundNewOwnerID).Return(true)

	handler := New(storage)

	params := &Params{OwnerID: GoodOwnerID, ItemID: GoodItemID, NewOwnerID: NotFoundNewOwnerID}
	result, err := handler.Execute(params)
	if assert.NoError(t, err) {
		assert.IsType(t, result, &Result{})
	}
}

func TestItemNotFound(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("ChangeUserItemOwner", GoodOwnerID, NotFoundItemID, GoodNewOwnerID).Return(true)

	handler := New(storage)

	params := &Params{OwnerID: GoodOwnerID, ItemID: NotFoundItemID, NewOwnerID: GoodNewOwnerID}
	result, err := handler.Execute(params)
	if assert.NoError(t, err) {
		assert.IsType(t, result, &Result{})
	}
}
