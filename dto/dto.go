package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	// ====================================== Failed ======================================
	MESSAGE_INVALID_REQUEST_PAYLOAD = "invalid request payload"

	// Middleware
	MESSAGE_FAILED_PROSES_REQUEST      = "failed proses request"
	MESSAGE_FAILED_ACCESS_DENIED       = "failed access denied"
	MESSAGE_FAILED_TOKEN_NOT_FOUND     = "failed token not found"
	MESSAGE_FAILED_TOKEN_NOT_VALID     = "failed token not valid"
	MESSAGE_FAILED_TOKEN_DENIED_ACCESS = "failed token denied access"
	MESSAGE_FAILED_GET_CUSTOM_CLAIMS   = "failed get custom claims"
	MESSAGE_FAILED_GET_ROLE_USER       = "failed get role user"

	// Query Params
	MESSAGE_INVALID_QUERY_PARAMS = "invalid query params"

	// UUID
	MESSAGE_FAILED_INVALID_UUID = "invalid UUID format"

	// File
	MESSAGE_FAILED_PARSE_MULTIPART_FORM = "failed to parse multipart form"
	MESSAGE_FAILED_NO_FILES_UPLOADED    = "failed no files uploaded"
	MESSAGE_FAILED_UPLOAD_FILES         = "failed upload files"

	// Authentication Errors
	FAILED_SIGNIN        = "failed signin"
	FAILED_REFRESH_TOKEN = "failed refresh token"

	// logging Errors
	FAILED_CREATE_LOG             = "failed to create log"
	FAILED_GET_LOGS               = "failed to get logs"
	FAILED_GET_LOGS_BY_STORE_ID   = "failed to get logs by store ID"
	FAILED_GET_LOGS_BY_DATE_RANGE = "failed to get logs by date range"

	// User Errors
	FAILED_GET_PROFILE = "failed to get profile"

	// General Errors
	FAILED_CREATE         = "failed to create"
	FAILED_UPDATE         = "failed to update"
	FAILED_DELETE         = "failed to delete"
	FAILED_GET_ALL        = "failed to get all"
	FAILED_GET_DETAIL     = "failed to get detail"
	NOT_FOUND             = "not found"
	INTERNAL_SERVER_ERROR = "internal server error"

	// ====================================== Success ======================================
	// File
	MESSAGE_SUCCESS_UPLOAD_FILES = "success upload files"
	MESSAGE_SUCCESS_UPLOAD_FILE  = "success upload file"

	// Authentication Sucess
	SUCCESS_SIGNIN        = "success signin"
	SUCCESS_REFRESH_TOKEN = "success refresh token"

	// Logging Success
	SUCCESS_CREATE_LOG             = "success to create log"
	SUCCESS_GET_LOGS               = "success to get logs"
	SUCCESS_GET_LOGS_BY_STORE_ID   = "success to get logs by store ID"
	SUCCESS_GET_LOGS_BY_DATE_RANGE = "success to get logs by date range"

	// User Errors
	SUCCESS_GET_PROFILE = "sucess to get profile"

	// General Success
	SUCCESS_CREATE     = "success create"
	SUCCESS_UPDATE     = "success update"
	SUCCESS_DELETE     = "success delete"
	SUCCESS_GET_ALL    = "success get all"
	SUCCESS_GET_DETAIL = "success get detail"
)

var (

	// Token
	ErrGenerateAccessToken           = errors.New("failed to generate access token")
	ErrGenerateRefreshToken          = errors.New("failed to generate refresh token")
	ErrGenerateAccessAndRefreshToken = errors.New("failed to generate access token and refresh token")
	ErrUnexpectedSigningMethod       = errors.New("unexpected signing method")
	ErrDecryptToken                  = errors.New("failed to decrypt token")
	ErrTokenInvalid                  = errors.New("token invalid")
	ErrValidateToken                 = errors.New("failed to validate token")
	ErrGetUserIDFromToken            = errors.New("failed get user id from token")

	// File
	ErrNoFilesUploaded    = errors.New("failed no files uploaded")
	ErrInvalidFileType    = errors.New("invalid file type")
	ErrSaveFile           = errors.New("failed save file")
	ErrCreateFolderAssets = errors.New("failed create folder assets")
	ErrDeleteOldImage     = errors.New("failed to delete old image")

	// General
	ErrNotFound         = errors.New("not found")
	ErrValidationFailed = errors.New("validation failed")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInternal         = errors.New("error internal")
	ErrUnauthorized     = errors.New("unauthorized")

	// Input

	// Authentication
	ErrIncorrectPassword = errors.New("credential incorrect")

	// Logging
	ErrCreateLog          = errors.New("failed to create log")
	ErrGetLogs            = errors.New("failed to get logs")
	ErrGetLogsByStoreID   = errors.New("failed to get logs by store ID")
	ErrGetLogsByDateRange = errors.New("failed to get logs by date range")

	// User
	ErrGetUserByEmail = errors.New("failed to get user by email")
	ErrGetUserByID    = errors.New("failed get user by id")

	// Parse
)

// Pagination
type (
	LoggingPaginationResponse struct {
		Data       []LoggingResponse  `json:"data"`
		Pagination PaginationResponse `json:"pagination"`
	}
)

// Without Pagination
type (
	// Authentication
	SignInRequest struct {
		Email    string `json:"email" binding:"required" example:"admin@mail.com"`
		Password string `json:"password" binding:"required" example:"secret123"`
	}
	SignInResponse struct {
		AccessToken  string `json:"access_token" example:"<access_token_here>"`
		RefreshToken string `json:"refresh_token" example:"<refresh_token_here>"`
	}
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required" example:"<refresh_token_here>"`
	}
	RefreshTokenResponse struct {
		AccessToken string `json:"access_token" binding:"required" example:"<new_access_token_here>"`
	}

	// Logging
	LoggingRequest struct {
		StoreID *uuid.UUID `json:"store_id,omitempty"`
		Action  string     `json:"action" example:"Create"`
		Message string     `json:"message" binding:"required" example:"Created a new store"`
	}
	LoggingResponse struct {
		ID        uuid.UUID  `json:"id"`
		StoreID   *uuid.UUID `json:"store_id,omitempty"`
		Action    string     `json:"action" example:"Create"`
		Message   string     `json:"message" example:"Created a new store"`
		CreatedAt time.Time  `json:"created_at"`
	}
)

// User
type (
	UserResponse struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
		Name  string    `json:"name"`
	}
)
