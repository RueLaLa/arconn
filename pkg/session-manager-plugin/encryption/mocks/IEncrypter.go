// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	log "github.com/ruelala/arconn/pkg/session-manager-plugin/log"
	mock "github.com/stretchr/testify/mock"
)

// IEncrypter is an autogenerated mock type for the IEncrypter type
type IEncrypter struct {
	mock.Mock
}

// Decrypt provides a mock function with given fields: _a0, cipherText
func (_m *IEncrypter) Decrypt(_a0 log.T, cipherText []byte) ([]byte, error) {
	ret := _m.Called(_a0, cipherText)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(log.T, []byte) []byte); ok {
		r0 = rf(_a0, cipherText)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(log.T, []byte) error); ok {
		r1 = rf(_a0, cipherText)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Encrypt provides a mock function with given fields: _a0, plainText
func (_m *IEncrypter) Encrypt(_a0 log.T, plainText []byte) ([]byte, error) {
	ret := _m.Called(_a0, plainText)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(log.T, []byte) []byte); ok {
		r0 = rf(_a0, plainText)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(log.T, []byte) error); ok {
		r1 = rf(_a0, plainText)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEncryptedDataKey provides a mock function with given fields:
func (_m *IEncrypter) GetEncryptedDataKey() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}
