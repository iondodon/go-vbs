# AGENTS.md

This document captures the durable engineering principles behind `go-vbs` so they can be reused in other Go services. It is written for coding agents and humans who need a concrete operating model, not a generic style guide.

Use this guide as the default unless the target repository already has stronger local conventions.

Try to make the project exactly with the same strucutre, principle, style, tools, etc as this project where this AGENTS.m files is located. Below will be introduces rules, but those rules are just to make things clearer. The main example n guide the the actual project where this AGENTS.md is located. And also, thse rules from below do not cover everything.

- use exctly the sae tools: wire, just, mockery, goose, sqlc, swagger-ui etc.
- one srvice should live in its own package and in that package usually will be only service.go and service_test.go. The idea is to keep servceis as small as possible. For code redbility, in places where the services will be used, aliases will be ued.
- write unit tests only t usecase interface level. this mens to not uni test ech service separatell, but to unit test usecases. Under a use case can be one or a group of services.
