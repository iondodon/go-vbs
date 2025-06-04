# go-vbs

This is [VBS](https://github.com/iondodon/vbs) (originally implemented in Java) project reimplemented in Go.

## Tools Used

- goose
- sqlc
- mockery
- swagger-ui - dist/ from [https://github.com/swagger-api/swagger-ui/releases](https://github.com/swagger-api/swagger-ui/releases)
- **Google Wire** - Dependency injection code generation

## Dependency Injection with Google Wire

This project uses [Google Wire](https://github.com/google/wire) for compile-time dependency injection. Wire generates code that wires up all dependencies automatically, eliminating the need for manual dependency construction.

### Wire Setup

**Installation:**

```bash
# Install the Wire command-line tool
go install github.com/google/wire/cmd/wire@latest

# Or use the Makefile
make install-wire
```

**Files:**

- `wire.go` - Contains provider functions and Wire configuration (build tag: `wireinject`)
- `wire_gen.go` - Generated code by Wire (build tag: `!wireinject`)
- `Makefile` - Includes Wire commands for convenience

### Provider Functions

Each component has a dedicated provider function in `wire.go`:

```go
// Logger providers with wrapper types to avoid conflicts
func ProvideInfoLogger() InfoLogger
func ProvideErrorLogger() ErrorLogger

// Repository layer providers
func ProvideQueries(db *sql.DB) *repository.Queries
func ProvideVehicleRepository(queries *repository.Queries) vehicleBusiness.VehicleRepository

// Business layer providers
func ProvideGetVehicleUseCase(vehicleRepo vehicleBusiness.VehicleRepository) vehicleBusiness.GetVehicleUseCase
func ProvideBookVehicleUseCase(...) business.BookVehicleUseCase

// Controller providers
func ProvideVehicleController(...) *vehicleController.Controller
```

### Provider Sets

Related providers are grouped into sets for better organization:

```go
var RepositorySet = wire.NewSet(
    ProvideQueries,
    ProvideVehicleRepository,
    ProvideCustomerRepository,
    ProvideBookingRepository,
    ProvideBookingDateRepository,
)

var VehicleSet = wire.NewSet(
    ProvideGetVehicleUseCase,
    ProvideAvailabilityUseCase,
    ProvideVehicleController,
)

var ApplicationSet = wire.NewSet(
    LoggerSet,
    RepositorySet,
    VehicleSet,
    BookingSet,
    AuthSet,
)
```

### Injector Functions

Wire generates these functions for dependency injection:

```go
func InitializeAuthController(db *sql.DB) (*authController.Controller, error)
func InitializeVehicleController(db *sql.DB) (*vehicleController.Controller, error)
func InitializeBookingController(db *sql.DB) (*bookingController.Controller, error)
```

### Usage in main.go

```go
func main() {
    // ... database setup ...

    // Initialize controllers using Wire
    authController, err := InitializeAuthController(db)
    if err != nil {
        errorLog.Fatal(err)
    }

    vehicleController, err := InitializeVehicleController(db)
    if err != nil {
        errorLog.Fatal(err)
    }

    bookingController, err := InitializeBookingController(db)
    if err != nil {
        errorLog.Fatal(err)
    }

    // Use controllers in HTTP routes...
}
```

### Development Workflow

**Regenerating Wire code:**

```bash
# Using Wire directly
wire

# Using Makefile
make wire

# Build (includes wire generation)
make build
```

**When to regenerate:**

- After adding new dependencies
- After modifying provider functions
- After changing provider sets

### Logger Handling

Wire requires unique types for dependency injection. Since multiple components need loggers, we use wrapper types:

```go
type InfoLogger struct {
    *log.Logger
}

type ErrorLogger struct {
    *log.Logger
}
```

This allows Wire to distinguish between different logger types while maintaining the same underlying `*log.Logger` interface.

### Cross-Domain Dependencies

Wire elegantly handles cross-domain dependencies through interfaces:

```go
// Booking domain needs vehicle functionality
func ProvideBookVehicleUseCase(
    vehicleRepo vehicleBusiness.VehicleRepository,  // From vehicle domain
    customerRepo customerBusiness.CustomerRepository, // From customer domain
    vehicleAvailabilityService vehicleBusiness.AvailabilityUseCase, // From vehicle domain
    // ... other dependencies
) business.BookVehicleUseCase
```

Wire automatically resolves these dependencies from the appropriate domain providers.

## API Testing

### REST Client Extension

This project includes a `requests.http` file for testing the API endpoints using the [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) VSCode extension.

**Setup:**

1. Open the project in VSCode
2. Install the recommended REST Client extension (VSCode will prompt you automatically)
3. Open the `requests.http` file
4. Start the application (`go run .`)
5. Click "Send Request" above any endpoint to test it

**Available endpoints:**

- `GET /login` - Get access token
- `GET /refresh` - Refresh token
- `GET /vehicles/{uuid}` - Get vehicle by UUID
- `POST /bookings` - Create a booking
- `GET /bookings` - Get all bookings (requires JWT)

**VSCode Configuration:**

- Extension recommendation in `.vscode/extensions.json`
- REST Client settings configured in `.vscode/settings.json`
- Environment variables pre-configured for localhost development

### Swagger UI (Development)

Set `GO_ENV=development` to enable Swagger UI at `http://localhost:8000/docs`

## Architecture & Design Rules

This project follows specific architectural patterns and naming conventions to maintain clean architecture principles and clear separation of concerns.

### 1. Domain-Based Organization

**The codebase is organized by domain rather than by technical layers.**

- ✅ **Domain packages**: Each domain (vehicle, booking, customer, auth) contains all its layers
- ✅ **Self-contained**: Each domain can potentially be extracted as a separate module
- ✅ **Clear boundaries**: Dependencies between domains are explicit and controlled
- ✅ **Core Entities**: Each domain defines its core entities as structs within a `domain` sub-package (e.g., `vehicle/domain/Vehicle.go`, `booking/domain/Booking.go`). These represent the fundamental concepts of that domain.
- ✅ **Data Transfer Objects (DTOs)**: DTOs are used for structured data exchange between layers, such as controller request payloads or service responses. DTOs are located in the `controller/` directory (e.g., `booking/controller/DatePeriodDTO.go`).

### 2. Domain Structure with Individual Package Organization

Each domain package follows this enhanced structure where each struct has its own package:

```go
domain/
├── business/
│   ├── usecases.go           # Use case interfaces (what external systems can call)
│   ├── repositories.go       # Repository interfaces (what business logic needs)
│   ├── structName/           # Individual packages for each service
│   │   └── service.go        # Service implementation (e.g., getVehicle.Service)
│   └── anotherStruct/
│       └── service.go        # Another service implementation
├── controller/               # Input adapters (receive from outside world)
│   └── structController/     # Individual controller packages
│       └── controller.go     # Controller implementation (e.g., vehicleController.Controller)
└── repository/               # Output adapters (send to outside world)
    └── structRepository/     # Individual repository packages
        └── repository.go     # Repository implementation (e.g., vehicleRepository.Repository)
```

**✅ Clean Architecture Compliance:**

- **Business layer** only sees domain entities and interfaces ✅
- **Input adapters (controller/)** depend on business interfaces ✅
- **Output adapters (repository/)** implement business interfaces ✅
- **Business logic** CANNOT see any adapters (controller/ or repository/) ✅
- **Dependencies flow inward** (controller/repository → business → domain) ✅

### 3. Package Naming Conventions

**Each struct gets its own package using camelCase naming:**

- ✅ **Service packages**: `getVehicleService`, `availabilityService`, `bookVehicleService`, `getAllBookingsService`
- ✅ **Repository packages**: `vehicleRepository`, `bookingRepository`, `customerRepository`
- ✅ **Controller packages**: `vehicleController`, `bookingController`, `authController`
- ✅ **Struct naming**: All structs are named `Service`, `Repository`, or `Controller` within their packages

### 4. Interface Usage Rules

**Interfaces are used ONLY for layer boundaries and cross-domain dependencies.**

- ✅ **Use interfaces for**: Repository → Business dependencies (infrastructure boundary)
- ✅ **Use interfaces for**: Business → Controller dependencies (application boundary)
- ✅ **Use interfaces for**: Cross-domain dependencies (domain boundaries)
- ❌ **Do NOT use interfaces for**: Internal domain dependencies (same domain)

### 5. Interface Boundaries (Ports & Adapters)

Each domain defines clear boundaries using separate files:

#### Use Case Interfaces (`usecases.go`)

- **Use case interfaces**: Define what external systems can call into the business logic
- **Cross-domain services**: Services from other domains that this domain uses
- **Examples**: `BookVehicleUseCase`, `GetVehicleUseCase`, `VehicleAvailabilityService`

#### Repository Interfaces (`repositories.go`)

- **Repository interfaces**: Define what the business logic needs from data access
- **Cross-domain repositories**: Data access interfaces from other domains
- **Examples**: `VehicleRepository`, `CustomerRepository`, `BookingRepository`

### 6. Cross-Domain Dependencies

Cross-domain dependencies are handled through interfaces defined in the consuming domain's interface files:

```go
// In booking/business/repositories.go - booking domain needs data from other domains
type VehicleRepository interface {
    FindByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
    VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}

// In booking/business/usecases.go - booking domain needs services from other domains
type VehicleAvailabilityService interface {
    CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}
```

### 7. Current Domain Structure

```
vehicle/
├── domain/
│   ├── vehicle.go                  # Vehicle domain entity
│   └── vehicleCategory.go          # VehicleCategory domain entity
├── business/
│   ├── usecases.go                 # GetVehicleUseCase, AvailabilityUseCase interfaces
│   ├── repositories.go             # VehicleRepository interface
│   ├── getVehicleService/
│   │   └── service.go             # getVehicleService.Service
│   └── availabilityService/
│       └── service.go             # availabilityService.Service
├── controller/
│   └── vehicleController/
│       └── controller.go          # vehicleController.Controller
└── repository/
    └── vehicleRepository/
        └── repository.go          # vehicleRepository.Repository

booking/
├── domain/
│   ├── booking.go                  # Booking domain entity
│   └── bookingDate.go             # BookingDate domain entity
├── controller/
│   ├── bookingController/
│   │   └── controller.go          # bookingController.Controller
│   ├── datePeriodDTO.go           # DatePeriodDTO for API requests
│   └── createBookingRequestDTO.go # CreateBookingRequestDTO for API requests
├── business/
│   ├── usecases.go                 # BookVehicleUseCase, GetAllBookingsUseCase, VehicleAvailabilityService interfaces
│   ├── repositories.go             # Repository interfaces + cross-domain repository interfaces
│   ├── bookVehicleService/
│   │   └── service.go             # bookVehicleService.Service
│   └── getAllBookingsService/
│       └── service.go             # getAllBookingsService.Service
└── repository/
    ├── bookingRepository/
    │   └── repository.go          # bookingRepository.Repository
    └── bookingDateRepository/
        └── repository.go          # bookingDateRepository.Repository

customer/
├── domain/
│   └── customer.go                 # Customer domain entity
├── business/
│   └── repositories.go             # CustomerRepository interface
└── repository/
    └── customerRepository/
        └── repository.go          # customerRepository.Repository

auth/
└── controller/
    └── authController/
        └── controller.go          # authController.Controller

handler/
└── handler.go                     # HTTP error handling wrapper

repository/
├── db.go                         # Database connection utilities
├── integration_db.go             # Database connection function (moved from integration)
├── models.go                     # Generated SQLC models
├── query.sql.go                  # Generated SQLC queries
├── query.sql                     # SQL queries for SQLC
└── migrations/                   # Database migrations (moved from integration)
    ├── 20250114092700_initialize_schema.sql
    └── 20250114092854_add_mock_data.sql
```

### 8. Adapter Layers (Ports & Adapters)

The `controller/` and `repository/` directories implement the **Ports & Adapters** pattern:

#### Input Adapters (`controller/`)

- **Purpose**: Receive input from the outside world and translate it to business operations
- **Examples**: HTTP controllers, CLI handlers, gRPC servers, message queue consumers
- **Dependencies**: Import and use business interfaces from `usecases.go`
- **Current**: HTTP controllers that handle REST API requests

#### Output Adapters (`repository/`)

- **Purpose**: Send output to the outside world as requested by business logic
- **Examples**: Database repositories, external API clients, file systems, message queues
- **Dependencies**: Implement business interfaces from `repositories.go`
- **Current**: Database repositories using SQLC-generated queries

**Benefits:**

- 🎯 **Clear Responsibility**: Input vs Output adapters have distinct purposes
- 🔄 **Easy Replacement**: Swap HTTP for gRPC, or PostgreSQL for MongoDB
- 🧪 **Testability**: Mock adapters easily at domain boundaries
- 📦 **Future Growth**: Add new adapter types (WebSocket, GraphQL, etc.)

### 9. Interface File Organization

**Clear separation of interface types:**

#### `usecases.go` - Service/Use Case Interfaces

- **Business use cases**: What the domain can do for external consumers
- **Cross-domain services**: Services from other domains that this domain uses
- **Examples**: `BookVehicleUseCase`, `GetAllBookingsUseCase`, `VehicleAvailabilityService`

#### `repositories.go` - Data Access Interfaces

- **Repository interfaces**: What the domain needs for data persistence
- **Cross-domain repositories**: Data access from other domains
- **Examples**: `BookingRepository`, `VehicleRepository`, `CustomerRepository`

**Benefits:**

- 🎯 **Logical Grouping**: Services vs Data Access are clearly separated
- 📁 **Easy Navigation**: Developers know exactly where to find interface types
- 🔍 **Quick Discovery**: Interface type is obvious from filename
- 📖 **Self-Documenting**: File names clearly indicate their purpose

### 10. Benefits of Individual Package Organization

- **Modularity**: Each struct is isolated in its own package for maximum modularity
- **Clear Ownership**: Each package has a single responsibility and clear purpose
- **Easy Navigation**: Find specific implementations quickly by package name
- **Import Clarity**: Explicit imports make dependencies crystal clear
- **Future Extraction**: Each package can easily become a separate module
- **Testability**: Easy to test and mock individual components

### 11. Dependency Injection

Dependencies are bootstrapped in `boot.go` with clear cross-domain dependency injection:

```go
func BootstrapApplication(db *sql.DB) *Dependencies {
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

    // Create repository layer (out adapters)
    queries := repository.New(db)

    // Vehicle domain
    var vehicleRepo vehicleBusiness.VehicleRepository = vehicleRepository.New(queries)
    var getVehicleUC vehicleBusiness.GetVehicleUseCase = getVehicleService.New(vehicleRepo)
    var vehicleAvailabilityService vehicleBusiness.AvailabilityUseCase = availabilityService.New(vehicleRepo)

    // Customer domain
    var customerRepo customerBusiness.CustomerRepository = customerRepository.New(queries)

    // Booking domain
    var bookingRepo business.BookingRepository = bookingRepository.New(queries)
    var bookingDateRepo business.BookingDateRepository = bookingDateRepository.New(queries)

    var bookVehicleUC business.BookVehicleUseCase = bookVehicleService.New(
        infoLog,
        errorLog,
        vehicleRepo,  // Cross-domain dependency (vehicle out implements booking business interface)
        customerRepo, // Cross-domain dependency (customer out implements booking business interface)
        bookingRepo,
        bookingDateRepo,
        vehicleAvailabilityService, // Cross-domain dependency
    )

    var getAllBookingsUC business.GetAllBookingsUseCase = getAllBookingsService.New(bookingRepo)

    // Create controller layer (in adapters)
    authCtrl := authController.New(infoLog, errorLog)
    vehicleCtrl := vehicleController.New(infoLog, errorLog, getVehicleUC)
    bookingCtrl := bookingController.New(infoLog, errorLog, db, bookVehicleUC, getAllBookingsUC)

    return &Dependencies{
        AuthController:    authCtrl,
        VehicleController: vehicleCtrl,
        BookingController: bookingCtrl,
    }
}
```

This enhanced architecture provides maximum modularity and clear separation while maintaining clean architecture principles and enabling easy extraction of components into separate modules or microservices when needed.
