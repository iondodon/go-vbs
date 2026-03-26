package vehicle

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	bookingController "github.com/iondodon/go-vbs/internal/booking/controller"
	bookingDomain "github.com/iondodon/go-vbs/internal/booking/domain"
	"github.com/iondodon/go-vbs/internal/booking/services/mocks"
	customerDomain "github.com/iondodon/go-vbs/internal/customer/domain"
	vehicleDomain "github.com/iondodon/go-vbs/internal/vehicle/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test data setup
func setupTestData() (uuid.UUID, uuid.UUID, bookingController.DatePeriodDTO, *vehicleDomain.Vehicle, *customerDomain.Customer) {
	customerUUID := uuid.New()
	vehicleUUID := uuid.New()
	period := bookingController.DatePeriodDTO{
		FromDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		ToDate:   time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC),
	}

	vehicle := &vehicleDomain.Vehicle{
		UUID:               vehicleUUID,
		RegistrationNumber: "ABC123",
		Make:               "Toyota",
		Model:              "Camry",
		FuelType:           vehicleDomain.Petrol,
	}

	customer := &customerDomain.Customer{
		UUID:     customerUUID,
		Username: "johndoe",
	}

	return customerUUID, vehicleUUID, period, vehicle, customer
}

func TestService_ForPeriod_Success(t *testing.T) {
	customerUUID, vehicleUUID, period, vehicle, customer := setupTestData()

	mockVehicleRepo := mocks.NewMockVehicleRepository(t)
	mockCustomerRepo := mocks.NewMockCustomerRepository(t)
	mockBookingRepo := mocks.NewMockBookingRepository(t)
	mockBookingDateRepo := mocks.NewMockBookingDateRepository(t)
	mockAvailabilityService := mocks.NewMockVehicleAvailabilityService(t)

	service := New(
		mockVehicleRepo,
		mockCustomerRepo,
		mockBookingRepo,
		mockBookingDateRepo,
		mockAvailabilityService,
	)

	ctx := context.Background()
	tx := &sql.Tx{}

	mockAvailabilityService.EXPECT().
		CheckForPeriod(ctx, vehicleUUID, period).
		Return(true, nil).
		Once()

	mockVehicleRepo.EXPECT().
		FindByUUID(ctx, vehicleUUID).
		Return(vehicle, nil).
		Once()

	mockBookingDateRepo.EXPECT().
		FindAllInPeriodInclusive(ctx, period.FromDate, period.ToDate).
		Return([]*bookingDomain.BookingDate{}, nil).
		Once()

	mockBookingDateRepo.EXPECT().
		Save(ctx, tx, mock.MatchedBy(func(bd *bookingDomain.BookingDate) bool { return true })).
		Return(nil).
		Maybe()

	mockCustomerRepo.EXPECT().
		FindByUUID(ctx, customerUUID).
		Return(customer, nil).
		Once()

	mockBookingRepo.EXPECT().
		Save(ctx, tx, mock.MatchedBy(func(b *bookingDomain.Booking) bool { return true })).
		Return(nil).
		Once()

	err := service.ForPeriod(ctx, tx, customerUUID, vehicleUUID, period)

	assert.NoError(t, err)
}

func TestService_ForPeriod_VehicleNotAvailable(t *testing.T) {
	customerUUID, vehicleUUID, period, _, _ := setupTestData()

	mockAvailabilityService := mocks.NewMockVehicleAvailabilityService(t)

	service := New(
		nil,
		nil,
		nil,
		nil,
		mockAvailabilityService,
	)

	ctx := context.Background()
	tx := &sql.Tx{}

	mockAvailabilityService.EXPECT().
		CheckForPeriod(ctx, vehicleUUID, period).
		Return(false, nil).
		Once()

	err := service.ForPeriod(ctx, tx, customerUUID, vehicleUUID, period)

	assert.Error(t, err)
	expectedError := fmt.Sprintf(alreadyHired, vehicleUUID)
	assert.Contains(t, err.Error(), expectedError)
}

func TestService_ForPeriod_AvailabilityCheckError(t *testing.T) {
	customerUUID, vehicleUUID, period, _, _ := setupTestData()

	mockAvailabilityService := mocks.NewMockVehicleAvailabilityService(t)

	service := New(
		nil,
		nil,
		nil,
		nil,
		mockAvailabilityService,
	)

	ctx := context.Background()
	tx := &sql.Tx{}

	expectedError := fmt.Errorf("database connection failed")

	mockAvailabilityService.EXPECT().
		CheckForPeriod(ctx, vehicleUUID, period).
		Return(false, expectedError).
		Once()

	err := service.ForPeriod(ctx, tx, customerUUID, vehicleUUID, period)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
