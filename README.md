# Go gRPC Microservices

A collection of gRPC-based microservices built in Go, following a modular monorepo architecture with shared Protocol Buffer contracts.

## Architecture

```text
go-grpc-microservices/
├── proto/
│   ├── go.mod
│   └── user/
│       ├── user.proto
│       ├── user.pb.go
│       └── user_grpc.pb.go
│
├── user-service/
│   ├── cmd/
│   ├── internal/
│   ├── go.mod
│   └── go.sum
│
└── api-gateway/
```

## Services

### User Service

A gRPC microservice that provides CRUD operations for users.

#### RPC Methods

* CreateUser
* GetUser
* UpdateUser
* DeleteUser

## Tech Stack

* Go
* gRPC
* Protocol Buffers
* PostgreSQL
* golang-migrate

## Shared Proto Module

The `proto` module contains Protocol Buffer definitions and generated Go stubs shared across services.

Example import:

```go
import pb "github.com/nazibul7/go-grpc-microservices/proto/user"
```

## Getting Started

### Clone Repository

```bash
git clone <repository-url>
cd go-grpc-microservices
```

### Generate Protobuf Code

```bash
cd proto

protoc \
  --proto_path=. \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  user/user.proto
```

### Run User Service

```bash
cd user-service

go mod tidy
go run ./cmd
```

## Project Goals

* Learn gRPC fundamentals
* Build production-style microservices in Go
* Implement API Gateway patterns
* Explore service-to-service communication
* Add observability and monitoring
* Implement authentication and authorization
* Learn containerization and orchestration

## Roadmap

* [x] Shared Proto Module
* [x] User Service
* [ ] API Gateway
* [ ] Order Service
* [ ] Authentication Service
* [ ] Service Discovery
* [ ] Docker Support
* [ ] Kubernetes Deployment
* [ ] Monitoring and Tracing

## License

MIT