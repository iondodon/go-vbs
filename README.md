# go-vbs

First read: [all project's code](ALL_CODE.txt) - this was obtained using:

```shell
tree -a -f -i \
| grep -v -E "/\.git/|vbs\.db|go\.sum|LICENSE$|/go-vbs(/|$)|/swagger-ui/|/\.vscode(/|$)" \
| while read file; do
  [ -f "$file" ] && echo "===== $file =====" && cat "$file"
done | xclip -selection clipboard
```

This project is intended to serve as an example of how to create other Go projects.

It acts as a reference for project structure, architecture, tooling, layering, implementation style, and more.

Anything not described here should be inferred from the project itself.
The project is the main source of truth, while this document serves only as clarification.

Whenever a new Go project needs to be created, it should follow the rules and patterns defined in this project.

The project must use the following tools: **just, wire, sqlc, mockery, and swagger-ui**.

The project structure should include **cmd**, **internal**, and **pkg** directories:

- **pkg** is required only if the project is a library or if there are components that need to be exported.
- **internal** should be divided into **app**, **http**, **repository**, and domain modules.
- In this example, we have an **http** module because the project uses HTTP, but if other technologies are used, additional modules may be introduced.

It is not mandatory to have domain modules as shown in this example. At the beginning, everything can be placed in a single module. As the project grows, it can be split into domain modules.

If domain modules are introduced, communication between them must happen through **use case interfaces**. These interfaces should be defined in the **service package** of each domain module.

Each domain module (or the single module in smaller applications) must contain subpackages such as:

- **controller**
- **service**
- **repository**
- **domain**

### Important rules for services

Each service must have its own package, even if it is very small or simple.

Example:

```
booking/service/bookings/all/get
```

Each service package must:

- Contain a single `Service` struct
- Provide a `New` constructor function
- Usually expose one exported function

Additional notes:

- Helper functions inside the package should be unexported.
- Multiple exported functions are allowed only if they provide alternative ways to call the same service (similar to method overloading).
- These functions should remain closely related to maintain high cohesion.

Each service must have **high cohesion**, which is why:

- Each service should be small
- Each service should live in its own package

### Function signature rules

Exported service functions should follow this order of arguments:

1. `context.Context` (if needed)
2. transaction (if needed)
3. business arguments

Ideally, the total number of arguments should not exceed **3**.
If business logic requires more arguments, group them into a struct (similar to the Command pattern).

### Naming and usage

The action performed by a service is defined by its **package path**, not by:

- the function name
- or the struct name (which is always `Service`)

Examples:

```
booking/service/bookings/all/get
booking/service/bookings/get/all
```

When importing services, always use clear aliases.
Inside services, dependencies on other services should also have meaningful names, for example:

```
getAllBookings
```

### Service dependencies

Service packages must not directly access the infrastructure layer.

Instead:

- Define interfaces (e.g., for database access) inside the service package (e.g., `repository.go`)
- Use these interfaces for dependency injection

### Architecture

The project must follow **hexagonal architecture**:

- The **domain package** is the core and must not depend on anything else

- **Service packages** depend only on:
  - domain packages
  - other service packages

- External dependencies (infrastructure, other modules) are accessed via interfaces defined in:
  - `usecases.go` inside the service package

- The **infrastructure layer** calls services through these use case interfaces

### Repository and database

- All migrations must be located in:

  ```
  root/repository/migrations
  ```

- Database models should be defined in `models.go`
- SQL queries must be defined in `query.sql`
- Code is generated using **sqlc**

The **repository package** in each domain module:

- Calls sqlc-generated functions
- Maps database models to domain models

### Controllers

The **controller package** should:

- Contain grouped controllers in subpackages
- Define DTOs, organized into:
  - `request`
  - `response`

Transactions are usually handled in controllers but can be managed elsewhere if needed.
