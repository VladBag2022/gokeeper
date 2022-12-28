// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	proto "github.com/VladBag2022/gokeeper/internal/proto"
	mock "github.com/stretchr/testify/mock"
)

// GRPCStore is an autogenerated mock type for the GRPCStore type
type GRPCStore struct {
	mock.Mock
}

// DeleteMeta provides a mock function with given fields: ctx, id
func (_m *GRPCStore) DeleteMeta(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSecret provides a mock function with given fields: ctx, id
func (_m *GRPCStore) DeleteSecret(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEncryptedKey provides a mock function with given fields: ctx, userID
func (_m *GRPCStore) GetEncryptedKey(ctx context.Context, userID int64) (*proto.ClientSecret, error) {
	ret := _m.Called(ctx, userID)

	var r0 *proto.ClientSecret
	if rf, ok := ret.Get(0).(func(context.Context, int64) *proto.ClientSecret); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ClientSecret)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSecrets provides a mock function with given fields: ctx, userID
func (_m *GRPCStore) GetSecrets(ctx context.Context, userID int64) (*proto.ClientSecrets, error) {
	ret := _m.Called(ctx, userID)

	var r0 *proto.ClientSecrets
	if rf, ok := ret.Get(0).(func(context.Context, int64) *proto.ClientSecrets); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ClientSecrets)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsUserMeta provides a mock function with given fields: ctx, userID, metaID
func (_m *GRPCStore) IsUserMeta(ctx context.Context, userID int64, metaID int64) (bool, error) {
	ret := _m.Called(ctx, userID, metaID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) bool); ok {
		r0 = rf(ctx, userID, metaID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, userID, metaID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsUserSecret provides a mock function with given fields: ctx, userID, secretID
func (_m *GRPCStore) IsUserSecret(ctx context.Context, userID int64, secretID int64) (bool, error) {
	ret := _m.Called(ctx, userID, secretID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) bool); ok {
		r0 = rf(ctx, userID, secretID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, userID, secretID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsUsernameAvailable provides a mock function with given fields: ctx, username
func (_m *GRPCStore) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	ret := _m.Called(ctx, username)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignIn provides a mock function with given fields: ctx, credentials
func (_m *GRPCStore) SignIn(ctx context.Context, credentials *proto.Credentials) (int64, error) {
	ret := _m.Called(ctx, credentials)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *proto.Credentials) int64); ok {
		r0 = rf(ctx, credentials)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.Credentials) error); ok {
		r1 = rf(ctx, credentials)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignUp provides a mock function with given fields: ctx, credentials
func (_m *GRPCStore) SignUp(ctx context.Context, credentials *proto.Credentials) (int64, error) {
	ret := _m.Called(ctx, credentials)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *proto.Credentials) int64); ok {
		r0 = rf(ctx, credentials)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.Credentials) error); ok {
		r1 = rf(ctx, credentials)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreMeta provides a mock function with given fields: ctx, secretID, meta
func (_m *GRPCStore) StoreMeta(ctx context.Context, secretID int64, meta *proto.Meta) (int64, error) {
	ret := _m.Called(ctx, secretID, meta)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64, *proto.Meta) int64); ok {
		r0 = rf(ctx, secretID, meta)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, *proto.Meta) error); ok {
		r1 = rf(ctx, secretID, meta)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreSecret provides a mock function with given fields: ctx, userID, secret
func (_m *GRPCStore) StoreSecret(ctx context.Context, userID int64, secret *proto.Secret) (int64, error) {
	ret := _m.Called(ctx, userID, secret)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64, *proto.Secret) int64); ok {
		r0 = rf(ctx, userID, secret)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, *proto.Secret) error); ok {
		r1 = rf(ctx, userID, secret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMeta provides a mock function with given fields: ctx, id, meta
func (_m *GRPCStore) UpdateMeta(ctx context.Context, id int64, meta *proto.Meta) error {
	ret := _m.Called(ctx, id, meta)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *proto.Meta) error); ok {
		r0 = rf(ctx, id, meta)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSecret provides a mock function with given fields: ctx, id, secret
func (_m *GRPCStore) UpdateSecret(ctx context.Context, id int64, secret *proto.Secret) error {
	ret := _m.Called(ctx, id, secret)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *proto.Secret) error); ok {
		r0 = rf(ctx, id, secret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewGRPCStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewGRPCStore creates a new instance of GRPCStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGRPCStore(t mockConstructorTestingTNewGRPCStore) *GRPCStore {
	mock := &GRPCStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
