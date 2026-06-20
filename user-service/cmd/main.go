package main

import (
	"log"

	"github.com/nazibul7/go-grpc-microservices/user-service/internal/db"
	"github.com/nazibul7/go-grpc-microservices/user-service/internal/server"
)

func main() {
	connStr := `postgres://postgres:password@localhost:5432/grpcdb?sslmode=disable`
	if err := db.RunMigrations(connStr); err != nil {
		log.Fatal(err)
	}

	db, err := db.NewDatabase(connStr)
	if err != nil {
		log.Fatalf("DB connection error: %v\n", err)
	}

	// Ensure the connection is closed when main exits
	defer db.Close()
	
	lis, err := server.Listen(":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := server.NewGRPCServer(db)

	log.Println("gRPC server running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
