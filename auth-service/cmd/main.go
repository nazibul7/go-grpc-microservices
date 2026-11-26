package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/db"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL is not set")
	}

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

	log.Println("gRPC auth service server running on :50052")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
