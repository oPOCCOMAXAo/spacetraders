## Build & Generate

- `make build` - builds to `bin/stbot`
- `make generate` - runs swagger gen, templ fmt+gen, and `go generate ./...`
- `make lint` - golangci-lint with auto-fix enabled (5m timeout)
- `make dev-watch-templ` - auto-regenerate templ on changes
- Always run `make generate` before `make build` when .templ files or swagger annotations change

## Dependency Injection

- Each package has a `Package` struct and a `MakePackage()` constructor
- `MakePackage` params: own `Config` (value) + other `*Package` pointers (dependencies)
- Package structs expose service fields: `Service`, `Client`, etc.
- Services implement lifecycle interfaces (`Shutdownable`, `Servable`) and are registered via `lifecycle.Service.RegisterService()`
- Strict dependency hierarchy defined in `pkg/app/run.go` - lower packages never depend on higher ones

## Config

- Env vars via caarlos0/env, prefixed per service: `LOGGER_`, `SERVER_`, `TASKS_`, `LIFECYCLE_`
- Config struct in each package, aggregated in `pkg/config/config.go`

## Views & API docs

- Views: Templ (.templ files)

## Linter

- golangci-lint with all linters enabled by default, specific disables in `.golangci.yml`
- Auto-fix enabled (`issues.fix: true`) - linting modifies files in place
- linter is also formatter
- `MakePackage()`, `Serve()`, `Descriptions()`, `Commands()`, `Register()` calls exempted from funlen/cyclop

## TODO

User can add todos during your work. Never just remove them, fix by doing what todo says. `make lint` highlights these todos.

## Code style

- Never put dangling functions in service package, each function must be: method or be function in \*utils package.
- Tests for Service must be grouped into suite with testify/suite.Suite.
