package request

import (
	"errors"
	"net/http"
	"strconv"
)

// ParseID extracts and validates an integer ID from a URL path parameter.
// Returns the ID and true on success, writes HTTP error and returns false on failure.

func ParseID(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")

	if idStr == "" {
		return 0, errors.New("missing id")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
