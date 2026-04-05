# Go Microservice Template

This is a template for rapidly building Go microservices following the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) and industry best practices.

## Technology Stack

*   **Language:** Go ^1.24
*   **HTTP Router:** [Chi](https://github.com/go-chi/chi) (Lightweight and idiomatic)
*   **Internal Communication:** gRPC & Protocol Buffers
*   **Database / ORM:** PostgreSQL & [GORM](https://gorm.io/)
*   **Configuration:** [Viper](https://github.com/spf13/viper) (.yaml files and Environment Variables)
*   **Observability/Logging:** [Zap](https://github.com/uber-go/zap) (Structured JSON logging)

## Architecture & Best Practices

This project is built with a strong focus on:
*   **Graceful Shutdown:** Intercepts `SIGTERM` and `SIGINT` signals to ensure that existing connections (both HTTP and gRPC) and database transactions finish properly before the process exits.
*   **Context Propagation:** Extensive use of `context.Context` from the entry point (Handlers) all the way down to data access layers (DB), allowing for proper timeout and cancellation controls.
*   **Dependency Injection (DI):** No global variables or `init()` functions are used. Everything is instantiated in `cmd/server/main.go` using constructors (`NewX()`).
*   **Unified Error Handling:** An agnostic domain-error encapsulation layer under `pkg/errors` that automatically maps to specific HTTP and gRPC status codes at the transport level.
*   **Health Checks for K8s:** Exposed endpoints (`/healthz` and `/ready`) explicitly designed for Kubernetes liveness and readiness probes.

## Directory Structure

```plaintext
.
├── api/             # API Contracts and Protocol Buffers definitions
│   └── proto/       # .proto files (e.g., user.proto) and their auto-generated code
├── cmd/
│   └── server/      # Entrypoint (main func) and application wiring (DI).
├── internal/        # (Private) Specific logic for this microservice
│   ├── config/      # Viper loading and struct definitions
│   ├── handler/     # HTTP endpoint routing (Chi) and gRPC bindings.
│   ├── repository/  # Database access layer (PostgreSQL, GORM definitions).
│   ├── server/      # Go wrappers to spawn HTTP and gRPC Listeners.
│   └── service/     # Core Business logic of the entities.
├── pkg/             # General-purpose code (exportable for reuse)
│   ├── errors/      # Status wrapper and application error types.
│   └── logger/      # Global system logger initialization via Zap.
├── Makefile         # Build automation scripts
└── config.yaml      # Initial environment baseline variables (ports, db, etc.)
```

## Quick Start Guide

To get started, ensure you have `go`, `make`, and `docker` installed (Docker is used for spinning up a local db and compiling protoc).

### 1. Boot up Dependencies
The project requires PostgreSQL to run. You can easily start a temporary local database instance using Docker:
```bash
docker run -d --rm --name go-template-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=mydb -p 5432:5432 postgres:15
```

### 2. Synchronize Modules & Compile Proto
```bash
# Sync go.mod libraries
make tidy

# Compile the auto-generated gRPC code based on api/proto/v1/user.proto
make proto
```
*(Note: The `make proto` command might attempt to interact with your local protoc binary. If you do not have it configured, there is a ready-to-use Docker snippet provided in the Makefile).*

### 3. Run the Application
Start the microservice:
```bash
make run
```
> This will start the gRPC server on port `9090` and the HTTP (REST) server on port `8080`. Ports are configurable via `config.yaml` or environment variables.

## Testing Endpoints

Once the application is running, open another terminal strictly to test how the Request travels through the Service, Handler, and the Repository:

*   **Check Health (REST):**
    ```bash
    curl -v http://localhost:8080/healthz
    ```
*   **Create a User (REST POST to JSON):**
    ```bash
    curl -X POST http://localhost:8080/users \
       -H "Content-Type: application/json" \
       -d '{"name": "Developer", "email": "dev@company.com"}'
    ```
*   **Get User (Substitute UUID):**
    ```bash
    curl http://localhost:8080/users/YOUR-GENERATED-UUID-HERE
    ```

## Next Steps (Building upon this Template)
1.  **Rename the Module:** Change the module path from `github.com/juanpblasi/go-template` to your actual repository path doing a global Search & Replace.
2.  **Mutate the Domain:** Swap the word `user` in `internal/service/`, `handler`, and `repository` to reflect the main entity this microservice will handle.
3.  **Add New gRPC Models:** Edit the `.proto` files located inside `api/proto/` and run `make proto` to generate new bindings.
