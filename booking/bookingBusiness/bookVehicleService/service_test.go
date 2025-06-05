package bookVehicleService

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/booking/bookingBusiness/mocks"
	"github.com/iondodon/go-vbs/booking/bookingController"
	"github.com/iondodon/go-vbs/booking/bookingDomain"
	"github.com/iondodon/go-vbs/customer/customerDomain"
	"github.com/iondodon/go-vbs/vehicle/vehicleDomain"
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
	// Arrange
	customerUUID, vehicleUUID, period, vehicle, customer := setupTestData()

	// Create all generated mocks from the mocks package
	mockVehicleRepo := mocks.NewMockVehicleRepository(t)
	mockCustomerRepo := mocks.NewMockCustomerRepository(t)
	mockBookingRepo := mocks.NewMockBookingRepository(t)
	mockBookingDateRepo := mocks.NewMockBookingDateRepository(t)
	mockAvailabilityService := mocks.NewMockVehicleAvailabilityService(t)

	// Create loggers
	infoLog := log.New(os.Stdout, "INFO: ", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR: ", log.LstdFlags)

	// Create service with mocked dependencies
	service := New(
		infoLog,
		errorLog,
		mockVehicleRepo,
		mockCustomerRepo,
		mockBookingRepo,
		mockBookingDateRepo,
		mockAvailabilityService,
	)

	ctx := context.Background()
	tx := &sql.Tx{}

	// Set up mock expectations using the generated mock's EXPECT() method
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

	// Allow any number of Save calls for BookingDate (implementation details may vary)
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

	// Act
	err := service.ForPeriod(ctx, tx, customerUUID, vehicleUUID, period)

	// Assert
	assert.NoError(t, err)
}

func TestService_ForPeriod_VehicleNotAvailable(t *testing.T) {
	// Arrange
	customerUUID, vehicleUUID, period, _, _ := setupTestData()

	mockAvailabilityService := mocks.NewMockVehicleAvailabilityService(t)

	// Only need availability service for this test
	service := New(
		log.New(os.Stdout, "INFO: ", log.LstdFlags),
		log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		nil, // not needed for this test
		nil, // not needed for this test
		nil, // not needed for this test
		nil, // not needed for this test
		mockAvailabilityService,
	)

	ctx := context.Background()
	tx := &sql.Tx{}

	// Set up mock expectation - vehicle is NOT available
	mockAvailabilityService.EXPECT().
		CheckForPeriod(ctx, vehicleUUID, period).
		Return(false, nil).
		Once()

	// Act
	err := service.ForPeriod(ctx, tx, customerUUID, vehicleUUID, period)

	// Assert
	assert.Error(t, err)
	expectedError := fmt.Sprintf(alreadyHired, vehicleUUID)
	assert.Contains(t, err.Error(), expectedError)
}

func TestService_ForPeriod_AvailabilityCheckError(t *testing.T) {
	// Arrange
	customerUUID, vehicleUUID, period, _, _ := setupTestData()

	mockAvailabilityService := mocks.NewMockVehicleAvailabilityService(t)

	service := New(
		log.New(os.Stdout, "INFO: ", log.LstdFlags),
		log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		nil,
		nil,
		nil,
		nil,
		mockAvailabilityService,
	)

	ctx := context.Background()
	tx := &sql.Tx{}

	expectedError := fmt.Errorf("database connection failed")

	// Set up mock expectation - availability check fails
	mockAvailabilityService.EXPECT().
		CheckForPeriod(ctx, vehicleUUID, period).
		Return(false, expectedError).
		Once()

	// Act
	err := service.ForPeriod(ctx, tx, customerUUID, vehicleUUID, period)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
