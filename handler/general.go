package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
)

func mapErrorStatus(err error) int {
	switch {
	case errors.Is(err, dto.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, dto.ErrValidationFailed):
		return http.StatusBadRequest
	case errors.Is(err, dto.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, dto.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func mapErrorMessage(err error) string {
	switch {
	case errors.Is(err, dto.ErrNotFound):
		return "data not found"
	case errors.Is(err, dto.ErrValidationFailed):
		if !errors.Is(err, fmt.Errorf("validation failed")) {
			return err.Error()
		}
		return "validation failed"
	case errors.Is(err, dto.ErrAlreadyExists):
		return "data already exists"
	case errors.Is(err, dto.ErrUnauthorized):
		return "unauthorized access"
	default:
		return "internal server error"
	}
}
