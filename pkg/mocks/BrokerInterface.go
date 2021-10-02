// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import amqp "github.com/streadway/amqp"
import mock "github.com/stretchr/testify/mock"
import proto "github.com/golang/protobuf/proto"

// BrokerInterface is an autogenerated mock type for the BrokerInterface type
type BrokerInterface struct {
	mock.Mock
}

// Ping provides a mock function with given fields:
func (_m *BrokerInterface) Ping() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publish provides a mock function with given fields: _a0, _a1, _a2
func (_m *BrokerInterface) Publish(_a0 string, _a1 proto.Message, _a2 amqp.Table) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, proto.Message, amqp.Table) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PublishJson provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *BrokerInterface) PublishJson(_a0 string, _a1 interface{}, _a2 amqp.Table, _a3 int64) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, amqp.Table, int64) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PublishProto provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *BrokerInterface) PublishProto(_a0 string, _a1 proto.Message, _a2 amqp.Table, _a3 int64) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, proto.Message, amqp.Table, int64) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterSubscriber provides a mock function with given fields: _a0, _a1
func (_m *BrokerInterface) RegisterSubscriber(_a0 string, _a1 interface{}) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetExchangeName provides a mock function with given fields: _a0
func (_m *BrokerInterface) SetExchangeName(_a0 string) {
	_m.Called(_a0)
}

// SetQueueOptsArgs provides a mock function with given fields: args
func (_m *BrokerInterface) SetQueueOptsArgs(args amqp.Table) {
	_m.Called(args)
}

// Subscribe provides a mock function with given fields: _a0
func (_m *BrokerInterface) Subscribe(_a0 chan bool) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(chan bool) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SuppressAutoCreateQueue provides a mock function with given fields: _a0
func (_m *BrokerInterface) SuppressAutoCreateQueue(_a0 bool) {
	_m.Called(_a0)
}
