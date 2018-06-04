package handler_list

import (
	"testing"

	"github.com/denperov/owm-task/libs/types"
	"github.com/denperov/owm-task/services/items-api/handler_list/mocks"
	"github.com/stretchr/testify/assert"
)

var GoodOwnerID = types.Login("good_login")
var NotFoundOwnerID = types.Login("not_found_login")
var GoodItem = &types.Item{}
var GoodItems = []*types.Item{GoodItem}
var NotFoundItems = []*types.Item{}

func TestOK(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("GetUserItems", GoodOwnerID).Return(GoodItems)

	handler := New(storage)

	params := &Params{OwnerID: GoodOwnerID}
	result, err := handler.Execute(params)
	if assert.NoError(t, err) {
		if assert.IsType(t, result, &Result{}) {
			assert.ElementsMatch(t, result.(*Result).Items, GoodItems)
		}
	}
}

func TestNotFound(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("GetUserItems", NotFoundOwnerID).Return(NotFoundItems)

	handler := New(storage)

	params := &Params{OwnerID: NotFoundOwnerID}
	result, err := handler.Execute(params)
	if assert.NoError(t, err) {
		if assert.IsType(t, result, &Result{}) {
			assert.Empty(t, result.(*Result).Items)
		}
	}
}
