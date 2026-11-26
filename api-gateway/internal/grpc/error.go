package grpc

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandleGRPCError converts gRPC status errors into appropriate HTTP status codes.
// gRPC and HTTP have different error models — this is the translation layer.
func HandleGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		http.Error(w, st.Message(), http.StatusBadRequest)

	case codes.Unauthenticated:
		http.Error(w, st.Message(), http.StatusUnauthorized)

	case codes.PermissionDenied:
		http.Error(w, st.Message(), http.StatusForbidden)

	case codes.NotFound:
		http.Error(w, st.Message(), http.StatusNotFound)

	case codes.AlreadyExists:
		http.Error(w, st.Message(), http.StatusConflict)

	case codes.Unavailable:
		http.Error(w, "service temporarily unavailable", http.StatusServiceUnavailable)

	case codes.DeadlineExceeded:
		http.Error(w, "request timed out", http.StatusGatewayTimeout)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
