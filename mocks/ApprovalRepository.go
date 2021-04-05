// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/odpf/guardian/domain"
	mock "github.com/stretchr/testify/mock"
)

// ApprovalRepository is an autogenerated mock type for the ApprovalRepository type
type ApprovalRepository struct {
	mock.Mock
}

// BulkInsert provides a mock function with given fields: _a0
func (_m *ApprovalRepository) BulkInsert(_a0 []*domain.Approval) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*domain.Approval) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}