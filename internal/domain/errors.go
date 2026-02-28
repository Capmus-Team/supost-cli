package domain

import "errors"

// Domain errors. Designed to map cleanly to HTTP status codes.
// NotFound → 404, Validation → 400, Unauthorized → 401, Conflict → 409.
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
)
