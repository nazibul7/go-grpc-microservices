package main

import (
	"log"

	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/db"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/server"
)

func main() {
	connStr := `postgres://auth_user:auth_password@localhost:5433/auth_db?sslmode=disable`

	if err := db.RunMigrations(connStr); err != nil {
		log.Fatal(err)
	}
	
	db, err := db.NewDatabase(connStr)
	if err != nil {
		log.Fatalf("DB connection error: %v\n", err)
	}

	// Ensure the connection is closed when main exits
	defer db.Close()

	lis, err := server.Listen(":50052")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := server.NewGRPCServer(db)

	log.Println("gRPC server running on :50052")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
