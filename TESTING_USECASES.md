# Testing Use Cases in Go-VBS

This guide demonstrates how to test use cases in clean architecture using the generated mocks.

## Generated Mocks

The project uses mockery to generate mocks in organized `mocks` directories within each domain package:

- `booking/business/mocks/` - All booking business mocks
  - `booking_repository_mock.go`
  - `booking_date_repository_mock.go`
  - `customer_repository_mock.go`
  - `vehicle_repository_mock.go`
  - `book_vehicle_usecase_mock.go`
  - `get_all_bookings_usecase_mock.go`
  - `vehicle_availability_service_mock.go`

### Mock Configuration

```yaml
# .mockery.yaml
all: False
dir: "{{.InterfaceDir}}/mocks" # Generates in mocks subdirectory
outpkg: "mocks" # Package name is "mocks"
mockname: "Mock{{.InterfaceName}}"
with-expecter: true

packages:
  github.com/iondodon/go-vbs/booking/business:
    config:
    interfaces:
      BookingRepository:
        config:
          filename: "booking_repository_mock.go"
      # ... other interfaces
```

### Available Generated Mocks

From `booking/business/mocks/`:

- `mocks.NewMockVehicleRepository(t)`
- `mocks.NewMockCustomerRepository(t)`
- `mocks.NewMockBookingRepository(t)`
- `mocks.NewMockBookingDateRepository(t)`
- `mocks.NewMockVehicleAvailabilityService(t)`
- `mocks.NewMockBookVehicleUseCase(t)`
- `mocks.NewMockGetAllBookingsUseCase(t)`

## Testing Strategy

### 1. Import the Mocks Package

```go
import (
    "github.com/iondodon/go-vbs/booking/business/mocks"
)
```

### 2. Use Generated Mocks

```go
// Using generated mocks from the mocks package
mockVehicleRepo := mocks.NewMockVehicleRepository(t)
mockVehicleRepo.EXPECT().
    FindByUUID(ctx, vehicleUUID).
    Return(vehicle, nil).
    Once()
```

## Complete Test Example

```go
package bookVehicleService

import (
    "context"
    "database/sql"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/iondodon/go-vbs/booking/business/mocks"
    "github.com/iondodon/go-vbs/booking/controller"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestService_ForPeriod_Success(t *testing.T) {
    // Arrange - Setup test data
    customerUUID := uuid.New()
    vehicleUUID := uuid.New()
    period := controller.DatePeriodDTO{
        FromDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
        ToDate:   time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC),
    }

    // Create mocks from the mocks package
    mockVehicleRepo := mocks.NewMockVehicleRepository(t)
    mockAvailabilityService := mocks.NewMockVehicleAvailabilityService(t)
    mockCustomerRepo := mocks.NewMockCustomerRepository(t)
    mockBookingRepo := mocks.NewMockBookingRepository(t)
    mockBookingDateRepo := mocks.NewMockBookingDateRepository(t)

    // Create service
    service := New(
        mockVehicleRepo,
        mockCustomerRepo,
        mockBookingRepo,
        mockBookingDateRepo,
        mockAvailabilityService,
    )

    // Setup expectations
    mockAvailabilityService.EXPECT().
        CheckForPeriod(context.Background(), vehicleUUID, period).
        Return(true, nil).
        Once()

    mockVehicleRepo.EXPECT().
        FindByUUID(context.Background(), vehicleUUID).
        Return(vehicle, nil).
        Once()

    // Act
    err := service.ForPeriod(context.Background(), &sql.Tx{}, customerUUID, vehicleUUID, period)

    // Assert
    assert.NoError(t, err)
}
```

## Test Patterns

### 1. Happy Path Testing

Test the successful execution of the use case:

```go
func TestService_Success(t *testing.T) {
    // Setup all mocks to return successful responses
    mockRepo := mocks.NewMockRepository(t)
    mockRepo.EXPECT().FindByUUID(mock.Anything, mock.Anything).Return(entity, nil)

    // Execute use case
    result, err := service.Execute(ctx, input)

    // Verify success
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 2. Error Handling Testing

Test how the use case handles various error conditions:

```go
func TestService_EntityNotFound(t *testing.T) {
    // Setup mock to return error
    expectedError := errors.New("entity not found")
    mockRepo := mocks.NewMockRepository(t)
    mockRepo.EXPECT().FindByUUID(mock.Anything, mock.Anything).Return(nil, expectedError)

    // Execute use case
    result, err := service.Execute(ctx, input)

    // Verify error handling
    assert.Error(t, err)
    assert.Equal(t, expectedError, err)
    assert.Nil(t, result)
}
```

### 3. Business Logic Testing

Test specific business rules:

```go
func TestService_VehicleNotAvailable(t *testing.T) {
    // Setup availability service to return false
    mockAvailabilityService := mocks.NewMockVehicleAvailabilityService(t)
    mockAvailabilityService.EXPECT().
        CheckForPeriod(ctx, vehicleUUID, period).
        Return(false, nil)

    // Execute use case
    err := service.ForPeriod(ctx, tx, customerUUID, vehicleUUID, period)

    // Verify business rule enforcement
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already taken")
}
```

### 4. Transaction Testing

Test transaction handling:

```go
func TestService_TransactionRollback(t *testing.T) {
    // Setup mocks to succeed initially, then fail
    mockRepo1 := mocks.NewMockRepository1(t)
    mockRepo2 := mocks.NewMockRepository2(t)

    mockRepo1.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(nil)
    mockRepo2.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("save failed"))

    // Execute use case
    err := service.Execute(ctx, tx, input)

    // Verify error propagation (transaction should be rolled back by caller)
    assert.Error(t, err)
}
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./booking/business/bookVehicleService/ -v

# Run with coverage
go test ./booking/business/bookVehicleService/ -cover

# Run benchmarks
go test ./booking/business/bookVehicleService/ -bench=.
```

## Directory Structure

The organized structure looks like this:

```
booking/
├── business/
│   ├── mocks/                          # All generated mocks
│   │   ├── booking_repository_mock.go
│   │   ├── customer_repository_mock.go
│   │   ├── vehicle_repository_mock.go
│   │   ├── book_vehicle_usecase_mock.go
│   │   └── ...
│   ├── bookVehicleService/
│   │   ├── service.go
│   │   └── service_test.go             # Imports from mocks package
│   ├── repositories.go
│   └── usecases.go
```

## Best Practices

### 1. Test Data Setup

Create helper functions for test data:

```go
func setupTestData() (uuid.UUID, *domain.Entity) {
    id := uuid.New()
    entity := &domain.Entity{
        UUID: id,
        Name: "Test Entity",
    }
    return id, entity
}
```

### 2. Mock Organization

- Use generated mocks from the `mocks` package
- All mocks are organized in dedicated directories
- Clean separation from business logic
- Easy to find and maintain

### 3. Test Isolation

- Each test should be independent
- Use fresh mocks for each test
- Don't share state between tests

### 4. Assertion Strategy

- Test the specific behavior, not implementation details
- Verify both success and error cases
- Check that the right methods are called with correct parameters

## Generating New Mocks

To generate mocks for new interfaces:

1. Add the interface to `.mockery.yaml`
2. Run `make mocks` or `~/go/bin/mockery`
3. Import from the `mocks` package in your tests

Example:

```bash
# Add interface to config, then regenerate
make mocks
```

This will create the appropriate mock files in the `mocks` directories with type-safe mock implementations.

## Benefits of Organized Mock Structure

1. **Clean Separation** - Mocks are separated from business logic
2. **Easy Discovery** - All mocks for a domain are in one place
3. **No Clutter** - Business packages stay clean
4. **Professional Structure** - Industry-standard organization
5. **Package Clarity** - Clear import paths (`mocks.NewMock*`)
6. **Maintainable** - Easy to update and manage mock files
