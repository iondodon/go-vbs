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

### 2. Domain Structure with Individual Package Organization

Each domain package follows this enhanced structure where each struct has its own package:

```go
domain/
├── business/
│   ├── interfaces.go      # Business interfaces (Repository, UseCase contracts)
│   ├── service.go         # Backwards compatibility exports
│   ├── structName/        # Individual packages for each service
│   │   └── service.go     # Service implementation (e.g., getVehicle.Service)
│   └── anotherStruct/
│       └── service.go     # Another service implementation
├── in/                    # Input adapters (receive from outside world)
│   └── structController/  # Individual controller packages
│       └── controller.go  # Controller implementation (e.g., vehicleController.Controller)
└── out/                   # Output adapters (send to outside world)
    └── structRepository/  # Individual repository packages
        └── repository.go  # Repository implementation (e.g., vehicleRepository.Repository)
```

**✅ Clean Architecture Compliance:**

- **Business layer** only sees domain entities and interfaces ✅
- **Input adapters (in/)** depend on business interfaces ✅
- **Output adapters (out/)** implement business interfaces ✅
- **Business logic** CANNOT see any adapters (in/ or out/) ✅
- **Dependencies flow inward** (in/out → business → domain) ✅

### 3. Package Naming Conventions

**Each struct gets its own package using camelCase naming:**

- ✅ **Service packages**: `getVehicle`, `availability`, `bookVehicle`, `getAllBookings`
- ✅ **Repository packages**: `vehicleRepository`, `bookingRepository`, `customerRepository`
- ✅ **Controller packages**: `vehicleController`, `bookingController`, `authController`
- ✅ **Struct naming**: All structs are named `Service`, `Repository`, or `Controller` within their packages

### 4. Interface Usage Rules

**Interfaces are used ONLY for layer boundaries and cross-domain dependencies.**

- ✅ **Use interfaces for**: Repository → Business dependencies (infrastructure boundary)
- ✅ **Use interfaces for**: Business → Controller dependencies (application boundary)
- ✅ **Use interfaces for**: Cross-domain dependencies (domain boundaries)
- ❌ **Do NOT use interfaces for**: Internal domain dependencies (same domain)

### 5. Cross-Domain Dependencies

Cross-domain dependencies are handled through interfaces defined in the consuming domain:

```go
// In booking/business/interfaces.go - booking domain needs vehicle functionality
type VehicleRepository interface {
    FindByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
}

type VehicleAvailabilityService interface {
    CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}
```

### 6. Current Domain Structure

```
vehicle/
├── business/
│   ├── interfaces.go               # Repository and UseCase interfaces
│   ├── service.go                  # Backwards compatibility
│   ├── getVehicle/
│   │   └── service.go             # getVehicle.Service
│   └── availability/
│       └── service.go             # availability.Service
├── in/
│   └── vehicleController/
│       └── controller.go          # vehicleController.Controller
└── out/
    └── vehicleRepository/
        └── repository.go          # vehicleRepository.Repository

booking/
├── business/
│   ├── interfaces.go               # Repository and UseCase interfaces + cross-domain interfaces
│   ├── service.go                  # Backwards compatibility
│   ├── bookVehicle/
│   │   └── service.go             # bookVehicle.Service
│   └── getAllBookings/
│       └── service.go             # getAllBookings.Service
├── in/
│   └── bookingController/
│       └── controller.go          # bookingController.Controller
└── out/
    ├── bookingRepository/
    │   └── repository.go          # bookingRepository.Repository
    └── bookingDateRepository/
        └── repository.go          # bookingDateRepository.Repository

customer/
├── business/
│   └── interfaces.go               # Repository interface
└── out/
    └── customerRepository/
        └── repository.go          # customerRepository.Repository

auth/
└── in/
    └── authController/
        └── controller.go          # authController.Controller

handler/
└── handler.go                     # HTTP error handling wrapper

repository/
├── db.go                         # Database connection utilities
├── models.go                     # Generated SQLC models
├── query.sql.go                  # Generated SQLC queries
└── query.sql                     # SQL queries for SQLC
```

### 7. Adapter Layers (Ports & Adapters)

The `in/` and `out/` directories implement the **Ports & Adapters** pattern:

#### Input Adapters (`in/`)

- **Purpose**: Receive input from the outside world and translate it to business operations
- **Examples**: HTTP controllers, CLI handlers, gRPC servers, message queue consumers
- **Dependencies**: Import and use business interfaces
- **Current**: HTTP controllers that handle REST API requests

#### Output Adapters (`out/`)

- **Purpose**: Send output to the outside world as requested by business logic
- **Examples**: Database repositories, external API clients, file systems, message queues
- **Dependencies**: Implement business interfaces
- **Current**: Database repositories using SQLC-generated queries

**Benefits:**

- 🎯 **Clear Responsibility**: Input vs Output adapters have distinct purposes
- 🔄 **Easy Replacement**: Swap HTTP for gRPC, or PostgreSQL for MongoDB
- 🧪 **Testability**: Mock adapters easily at domain boundaries
- 📦 **Future Growth**: Add new adapter types (WebSocket, GraphQL, etc.)

### 8. Benefits of Individual Package Organization

- **Modularity**: Each struct is isolated in its own package for maximum modularity
- **Clear Ownership**: Each package has a single responsibility and clear purpose
- **Easy Navigation**: Find specific implementations quickly by package name
- **Import Clarity**: Explicit imports make dependencies crystal clear
- **Future Extraction**: Each package can easily become a separate module
- **Testability**: Easy to test and mock individual components

### 9. Benefits of Domain-Based Architecture

- **Modularity**: Each domain is self-contained and can be extracted as a separate module
- **Clear Ownership**: All code related to a domain is in one place
- **Reduced Coupling**: Dependencies between domains are explicit and minimal
- **Team Scalability**: Different teams can own different domains
- **Microservices Ready**: Easy to extract domains into separate services
- **Testability**: Easy to test entire domains in isolation

### 10. Dependency Injection

Dependencies are bootstrapped in `boot.go` with clear cross-domain dependency injection using specific packages:

```go
// Create output adapters (repositories) - using specific packages
vehicleRepo := vehicleRepository.NewRepository(queries)
customerRepo := customerRepository.NewRepository(queries)
bookingRepo := bookingRepository.NewRepository(queries)
bookingDateRepo := bookingDateRepository.NewRepository(queries)

// Create business services - using specific packages
getVehicleUC := getVehicle.NewService(vehicleRepo)
vehicleAvailabilityService := availability.NewService(vehicleRepo)

// Cross-domain dependency injection
bookVehicleUC := bookVehicle.NewService(
    infoLog,
    errorLog,
    vehicleRepo,                // Cross-domain dependency (vehicle out → booking business)
    customerRepo,               // Cross-domain dependency (customer out → booking business)
    bookingRepo,                // Same-domain dependency (booking out → booking business)
    bookingDateRepo,            // Same-domain dependency
    vehicleAvailabilityService, // Cross-domain dependency (vehicle business → booking business)
)

getAllBookingsUC := getAllBookings.NewService(bookingRepo)

// Create input adapters (controllers) - using specific packages
authCtrl := authController.NewController(infoLog, errorLog)
vehicleCtrl := vehicleController.NewController(infoLog, errorLog, getVehicleUC)
bookingCtrl := bookingController.NewController(infoLog, errorLog, db, bookVehicleUC, getAllBookingsUC)
```

This enhanced architecture provides maximum modularity and clear separation while maintaining clean architecture principles and enabling easy extraction of components into separate modules or microservices when needed.
