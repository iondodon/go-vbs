// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery

package vehicle

import (
	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
	mock "github.com/stretchr/testify/mock"
)

// NewMockVehicleRepository creates a new instance of MockVehicleRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockVehicleRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockVehicleRepository {
	mock := &MockVehicleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockVehicleRepository is an autogenerated mock type for the VehicleRepository type
type MockVehicleRepository struct {
	mock.Mock
}

type MockVehicleRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockVehicleRepository) EXPECT() *MockVehicleRepository_Expecter {
	return &MockVehicleRepository_Expecter{mock: &_m.Mock}
}

// FindByUUID provides a mock function for the type MockVehicleRepository
func (_mock *MockVehicleRepository) FindByUUID(vUUID uuid.UUID) (*domain.Vehicle, error) {
	ret := _mock.Called(vUUID)

	if len(ret) == 0 {
		panic("no return value specified for FindByUUID")
	}

	var r0 *domain.Vehicle
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) (*domain.Vehicle, error)); ok {
		return returnFunc(vUUID)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) *domain.Vehicle); ok {
		r0 = returnFunc(vUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Vehicle)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = returnFunc(vUUID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockVehicleRepository_FindByUUID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByUUID'
type MockVehicleRepository_FindByUUID_Call struct {
	*mock.Call
}

// FindByUUID is a helper method to define mock.On call
//   - vUUID
func (_e *MockVehicleRepository_Expecter) FindByUUID(vUUID interface{}) *MockVehicleRepository_FindByUUID_Call {
	return &MockVehicleRepository_FindByUUID_Call{Call: _e.mock.On("FindByUUID", vUUID)}
}

func (_c *MockVehicleRepository_FindByUUID_Call) Run(run func(vUUID uuid.UUID)) *MockVehicleRepository_FindByUUID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *MockVehicleRepository_FindByUUID_Call) Return(vehicle *domain.Vehicle, err error) *MockVehicleRepository_FindByUUID_Call {
	_c.Call.Return(vehicle, err)
	return _c
}

func (_c *MockVehicleRepository_FindByUUID_Call) RunAndReturn(run func(vUUID uuid.UUID) (*domain.Vehicle, error)) *MockVehicleRepository_FindByUUID_Call {
	_c.Call.Return(run)
	return _c
}

// VehicleHasBookedDatesOnPeriod provides a mock function for the type MockVehicleRepository
func (_mock *MockVehicleRepository) VehicleHasBookedDatesOnPeriod(vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error) {
	ret := _mock.Called(vUUID, period)

	if len(ret) == 0 {
		panic("no return value specified for VehicleHasBookedDatesOnPeriod")
	}

	var r0 bool
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID, dto.DatePeriodDTO) (bool, error)); ok {
		return returnFunc(vUUID, period)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID, dto.DatePeriodDTO) bool); ok {
		r0 = returnFunc(vUUID, period)
	} else {
		r0 = ret.Get(0).(bool)
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID, dto.DatePeriodDTO) error); ok {
		r1 = returnFunc(vUUID, period)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VehicleHasBookedDatesOnPeriod'
type MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call struct {
	*mock.Call
}

// VehicleHasBookedDatesOnPeriod is a helper method to define mock.On call
//   - vUUID
//   - period
func (_e *MockVehicleRepository_Expecter) VehicleHasBookedDatesOnPeriod(vUUID interface{}, period interface{}) *MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call {
	return &MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call{Call: _e.mock.On("VehicleHasBookedDatesOnPeriod", vUUID, period)}
}

func (_c *MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call) Run(run func(vUUID uuid.UUID, period dto.DatePeriodDTO)) *MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(dto.DatePeriodDTO))
	})
	return _c
}

func (_c *MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call) Return(b bool, err error) *MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call {
	_c.Call.Return(b, err)
	return _c
}

func (_c *MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call) RunAndReturn(run func(vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)) *MockVehicleRepository_VehicleHasBookedDatesOnPeriod_Call {
	_c.Call.Return(run)
	return _c
}
