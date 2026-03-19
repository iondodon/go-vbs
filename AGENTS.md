# AGENTS.md

This document captures the durable engineering principles behind `go-vbs` so they can be reused in other Go services. It is written for coding agents and humans who need a concrete operating model, not a generic style guide.

Use this guide as the default unless the target repository already has stronger local conventions.

Try to make the project exactly with the same structure, principles, style, tools, etc. as this project where this AGENTS.md file is located. Below will be introduced rules, but those rules are just to make things clearer. The main example and guide is the actual project where this AGENTS.md is located. Also, these rules from below do not cover everything.

- in `internal` package, besides other packages, there are domain `packages` which we can call modules. Those modules should have least coupling between them. In those modules will be put areas of the application that are not really related to one another. But, it is not a requirement for the application to be split into modules - only if it makes sense.

- each module has its `controller` package, `domain` package, `repository` package, `services` package and anything else that is needed - in other words - each module has its own layer in the hexagonal/onion architecture.

- keep all units (service, repository, domain struct, controller) as small as possible to keep the single responsibility principle and to provide more flexibility.

- one package per unit (service, controller, repository, etc). Not really necessary for domain because they are really simple and they can be grouped more in one single package, but still should be grouped in subpackages.

- each service should usually consist of only one single exported function, plus maybe a few unexported helper functions. There can be multiple exported functions, but those are only with purposes to provide the client of the service a more flexible API (like function overloading) by providing the possibility to call the service with different input parameters (a more flexible API). This would mean that those exported functions most probably will call each other internally in the service's package. One package per service.

- If context.Context is needed then it should be the first argument in the function. If a transaction is needed, then it should be after context, if context exists, otherwise, it should be the first.

- use exactly the same tools: wire, just, mockery, goose, sqlc, swagger-ui, etc.

- Swagger-ui should be downloaded from link. See `justfile` as a code example.

- one service should live in its own package and in that package usually will be only service.go and service_test.go. The idea is to keep services as small as possible. For code readability, in places where the services will be used, aliases will be used.

- write unit tests only at use case interface level. this means to not unit test each service separately, but to unit test use cases. Under a use case can be one or a group of services.
