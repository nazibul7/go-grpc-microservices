package router

import (
	"net/http"

	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/handler"
	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/middleware"
)

func RegisterRouter(mux *http.ServeMux, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) {
	mux.HandleFunc("POST /auth/signup", authHandler.SignUp)
	mux.HandleFunc("POST /auth/signin", authHandler.SignIn)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("POST /auth/signout", authHandler.SignOut)

	mux.Handle(
		"GET /user/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(userHandler.GetUser)),
	)

	mux.Handle(
		"PATCH /user/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(userHandler.UpdateUser)),
	)

	mux.Handle(
		"DELETE /user/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(userHandler.DeleteUser)),
	)
}
