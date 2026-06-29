package request

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// ParseID extracts and validates a UUID from the URL path.

func ParseID(r *http.Request) (string, error) {
	id := r.PathValue("id")

	if id == "" {
		return "", errors.New("missing id")
	}

	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("invalid id")
	}

	return id, nil
}
