package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
)

func mapErrorStatus(err error) int {
	switch {
	case errors.Is(err, dto.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, dto.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, dto.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, dto.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func cleanErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()
	parts := strings.Split(msg, ":")
	return strings.TrimSpace(parts[0])
}
