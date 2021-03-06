// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import types "github.com/denperov/owm-task/libs/types"

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// DeleteUserItem provides a mock function with given fields: ownerID, itemID
func (_m *Storage) DeleteUserItem(ownerID types.Login, itemID types.ItemID) bool {
	ret := _m.Called(ownerID, itemID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Login, types.ItemID) bool); ok {
		r0 = rf(ownerID, itemID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
