# go-vbs

This is [VBS](https://github.com/iondodon/vbs) (originally implemented in Java) project reimplemented in Go.

## Tools Used

- goose
- sqlc
- mockery
- swagger-ui - dist/ from [https://github.com/swagger-api/swagger-ui/releases](https://github.com/swagger-api/swagger-ui/releases)

## API Testing

### REST Client Extension

This project includes a `requests.http` file for testing the API endpoints using the [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) VSCode extension.

**Setup:**

1. Install the REST Client extension in VSCode
2. Open the `requests.http` file
3. Start the application (`go run .`)
4. Click "Send Request" above any endpoint to test it

**Available endpoints:**

- `GET /login` - Get access token
- `GET /refresh` - Refresh token
- `GET /vehicles/{uuid}` - Get vehicle by UUID
- `POST /bookings` - Create a booking
- `GET /bookings` - Get all bookings (requires JWT)

### Swagger UI (Development)

Set `GO_ENV=development` to enable Swagger UI at `http://localhost:8000/docs`

## Architecture & Design Rules

This project follows specific architectural patterns and naming conventions to maintain clean architecture principles and clear separation of concerns.

### 1. Interface Usage Rules

**Interfaces are used ONLY for layer boundaries, not for internal dependencies within the same layer.**

- ✅ **Use interfaces for**: Repository → Business dependencies (infrastructure boundary)
- ✅ **Use interfaces for**: Business → Controller dependencies (infrastructure boundary)
- ❌ **Do NOT use interfaces for**: Business → Business dependencies (internal dependencies)
- ❌ **Do NOT use interfaces for**: Repository → Repository dependencies (internal dependencies)

### 2. Naming Conventions

#### Business Layer

- **Interface name**: `UseCase`
- **Struct name**: `Service`
- **Example**:

  ```go
  type UseCase interface {
      Execute(ctx context.Context) error
  }

  type Service struct {
      // dependencies
  }
  ```

#### Repository Layer

- **Interface name**: `XRepository` (without "Interface" suffix)
- **Struct name**: `Repository`
- **Example**:

  ```go
  type VehicleRepository interface {
      FindByUUID(ctx context.Context, uuid UUID) (*Vehicle, error)
  }

  type Repository struct {
      // dependencies
  }
  ```

### 3. Constructor Patterns

**All `New` functions return concrete types, not interfaces.**

- Constructors return `*Service` or `*Repository` (concrete types)
- Callers use explicit interface assignment when needed

#### Business Constructors

```go
// Constructor returns concrete type
func New(repo business.VehicleRepository) *Service {
    return &Service{repo: repo}
}

// Usage with explicit interface assignment
var getVehicleUC getVehicle.UseCase = getVehicle.New(vehicleRepo)
```

#### Repository Constructors

```go
// Constructor returns concrete type
func New(queries *repository.Queries) *Repository {
    return &Repository{queries: queries}
}

// Usage with explicit interface assignment
var vehicleRepo business.VehicleRepository = vehicleRepository.New(queries)
```

### 4. Dependency Injection Patterns

#### Infrastructure Boundaries (Use Interfaces)

```go
// Repository to Business (infrastructure boundary)
var vehicleRepo business.VehicleRepository = vehicleRepository.New(queries)
getVehicleUC := getVehicle.New(vehicleRepo)

// Business to Controller (infrastructure boundary)
var bookVehicleUC bookVehicle.UseCase = bookVehicle.New(...)
controller := bookingController.New(..., bookVehicleUC)
```

#### Internal Dependencies (Use Concrete Types)

```go
// Business to Business dependencies (internal)
isAvailableUC := isVehicleAvailable.New(vehicleRepo)  // concrete type
getBookingDatesUC := getBookingDate.New(bookingDateRepo)  // concrete type

bookVehicleUC := bookVehicle.New(
    logger,
    vehicleRepo,
    customerRepo,
    bookingRepo,
    isAvailableUC,      // concrete type
    getBookingDatesUC,  // concrete type
)
```

### 5. Layer Structure

```
Controllers (Infrastructure)
    ↓ (interfaces)
Business Logic
    ↓ (interfaces)
Repositories (Infrastructure)
    ↓
External Data Sources
```

### 6. Benefits of This Architecture

- **Clear Boundaries**: Interfaces only where truly needed for layer separation
- **Reduced Complexity**: No unnecessary abstractions for internal dependencies
- **Explicit Dependencies**: Clear distinction between interface and concrete usage
- **Testability**: Easy to mock at layer boundaries while keeping internal logic simple
- **Maintainability**: Consistent patterns across the entire codebase

### 7. File Organization

```
business/
├── repository_interfaces.go          # Repository interfaces (layer boundary)
├── vehicle/
│   ├── getVehicle/
│   │   ├── usecase.go                # UseCase interface
│   │   └── service.go                # Service struct implementation
│   └── isVehicleAvailable/
│       └── service.go                # Service struct (no interface for internal use)
└── booking/
    ├── bookVehicle/
    │   ├── usecase.go                # UseCase interface
    │   └── service.go                # Service struct implementation
    └── getAllBookings/
        ├── usecase.go                # UseCase interface
        └── service.go                # Service struct implementation

repository/
├── vehicleRepository/
│   └── repository.go                 # Repository struct
└── bookingRepository/
    └── repository.go                 # Repository struct
```

This architecture ensures clean separation of concerns while avoiding over-engineering with unnecessary interfaces.
