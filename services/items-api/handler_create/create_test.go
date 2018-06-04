package handler_create

import (
	"testing"

	"github.com/denperov/owm-task/libs/types"
	"github.com/denperov/owm-task/services/items-api/handler_create/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var GoodOwnerID = types.Login("good_login")
var GoodContent = "good content"

func TestOK(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("CreateUserItem", GoodOwnerID, mock.AnythingOfType("*types.Item"))

	handler := New(storage)

	params := &Params{OwnerID: GoodOwnerID, Content: GoodContent}
	result, err := handler.Execute(params)
	if assert.NoError(t, err) {
		if assert.IsType(t, result, &Result{}) {
			assert.True(t, result.(*Result).ItemID.Validate())
		}
	}
}
