package main

import (
	"log"
	"net/http"

	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/client"
	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/handler"
)

func main() {
	conn, userClient, err := client.NewUserClient()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userHandler := handler.NewUserHandler(userClient)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /user", userHandler.CreateUser)
	mux.HandleFunc("GET /user/{id}", userHandler.GetUser)
	mux.HandleFunc("PATCH /user/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /user/{id}", userHandler.DeleteUser)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("API Gateway running on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
