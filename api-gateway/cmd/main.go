package main

import (
	"log"
	"net/http"

	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/client"
	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/handler"
	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/router"
)

func main() {
	conn, userClient, err := client.NewUserClient()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn, authClient, err := client.NewAuthClient()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userHandler := handler.NewUserHandler(userClient)
	authHandler := handler.NewAuthHandler(authClient)

	mux := http.NewServeMux()

	router.RegisterRouter(mux, authHandler, userHandler)
	
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("API Gateway running on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
