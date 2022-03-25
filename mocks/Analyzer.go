// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	database "mermerd/database"

	mock "github.com/stretchr/testify/mock"
)

// Analyzer is an autogenerated mock type for the Analyzer type
type Analyzer struct {
	mock.Mock
}

// Analyze provides a mock function with given fields:
func (_m *Analyzer) Analyze() (*database.Result, error) {
	ret := _m.Called()

	var r0 *database.Result
	if rf, ok := ret.Get(0).(func() *database.Result); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*database.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
