package router

import (
	"net/http"

	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/handler"
)

func RegisterRouter(mux *http.ServeMux, authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler) {
	mux.HandleFunc("POST /auth/signup", authHandler.SignUp)
	mux.HandleFunc("POST /auth/signin", authHandler.SignIn)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("POST /auth/signout", authHandler.SignOut)

	mux.HandleFunc("GET /user/{id}", userHandler.GetUser)
	mux.HandleFunc("PATCH /user/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /user/{id}", userHandler.DeleteUser)
}
