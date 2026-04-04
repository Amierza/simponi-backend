package dto

import (
	"errors"
	"time"

	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/response"
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
	MESSAGE_FAILED_TOKEN_EXPIRED       = "failed token expired"
	MESSAGE_FAILED_TOKEN_DENIED_ACCESS = "failed token denied access"
	MESSAGE_FAILED_GET_ROLE_USER       = "failed get role user"
	MESSAGE_FAILED_CHECK_PERMISSION    = "failed to check permission"
	MESSAGE_FAILED_FORBIDDEN           = "forbidden"

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

	// Product Errors
	FAILED_CREATE_PRODUCT           = "failed to create product"
	FAILED_UPDATE_PRODUCT           = "failed to update product"
	FAILED_DELETE_PRODUCT           = "failed to delete product"
	FAILED_GET_ALL_PRODUCTS         = "failed to get all products"
	FAILED_GET_PRODUCT_DETAIL       = "failed to get product detail"
	FAILED_GET_PRODUCTS_BY_CATEGORY = "failed to get products by category"
	FAILED_UPDATE_STOCK             = "failed to update stock"

	// General Errors
	FAILED_CREATE         = "failed to create"
	FAILED_UPDATE         = "failed to update"
	FAILED_DELETE         = "failed to delete"
	FAILED_GET_ALL        = "failed to get all"
	FAILED_GET_DETAIL     = "failed to get detail"
	NOT_FOUND             = "not found"
	INTERNAL_SERVER_ERROR = "internal server error"
	UNAUTHORIZED          = "unauthorized"

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

	// Product Errors
	SUCCESS_CREATE_PRODUCT           = "success create product"
	SUCCESS_UPDATE_PRODUCT           = "success update product"
	SUCCESS_DELETE_PRODUCT           = "success delete product"
	SUCCESS_GET_ALL_PRODUCTS         = "success get all products"
	SUCCESS_GET_PRODUCT_DETAIL       = "success get product detail"
	SUCCESS_GET_PRODUCTS_BY_CATEGORY = "success get products by category"
	SUCCESS_UPDATE_STOCK             = "success update stock"

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
	ErrNotFound      = errors.New("not found")
	ErrBadRequest    = errors.New("bad request")
	ErrAlreadyExists = errors.New("already exists")
	ErrInternal      = errors.New("error internal")
	ErrUnauthorized  = errors.New("unauthorized")

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

	// Product
	ErrCreateProduct           = errors.New("failed to create product")
	ErrUpdateProduct           = errors.New("failed to update product")
	ErrDeleteProduct           = errors.New("failed to delete product")
	ErrGetAllProducts          = errors.New("failed to get all products")
	ErrGetProductByID          = errors.New("failed to get product by id")
	ErrGetProductBySKU         = errors.New("failed to get product by sku")
	ErrGetProductsByCategory   = errors.New("failed to get products by category")
	ErrUpdateStock             = errors.New("failed to update stock")
	ErrProductSKUAlreadyExists = errors.New("product SKU already exists")

	// Parse
)

// Pagination
type (
	LogPaginationResponse struct {
		Data       []LogResponse               `json:"data"`
		Pagination response.PaginationResponse `json:"pagination"`
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

	// Log
	LogRequest struct {
		StoreID *uuid.UUID `json:"store_id,omitempty"`
		Action  string     `json:"action" example:"Create"`
		Message string     `json:"message" binding:"required" example:"Created a new store"`
	}
	LogResponse struct {
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

// Product

type (
	// Product Category
	ProductCategoryRequest struct {
		Name string `json:"name" binding:"required" example:"Electronics"`
	}
	ProductCategoryResponse struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name" example:"Electronics"`
		CreatedAt time.Time `json:"created_at"`
	}
	ProductCategoryPaginationResponse struct {
		response.PaginationResponse
		Data []ProductPaginationResponse `json:"data"`
	}
	ProductCategoryPaginationRepositoryResponse struct {
		response.PaginationResponse
		ProductCategories []entity.ProductCategory
	}
)

type (
	// Product Image
	ProductImageResponse struct {
		ID       uuid.UUID `json:"id"`
		ImageURL string    `json:"image_url" example:"https://example.com/image.jpg"`
	}
)

type (
	// External Product
	ExternalProductResponse struct {
		ID                uuid.UUID  `json:"id"`
		ProductID         *uuid.UUID `json:"product_id,omitempty"`
		StorePlatformID   *uuid.UUID `json:"store_platform_id,omitempty"`
		Price             int64      `json:"price" example:"150000"`
	}
)

type (
	// Product Stats
	ProductStatsResponse struct {
		TotalProducts int64 `json:"total_products"`
		TotalSKUs     int64 `json:"total_skus"`
		StockUnits    int64 `json:"stock_units"`
		LowStock      int64 `json:"low_stock"`
		OutOfStock    int64 `json:"out_of_stock"`
		Unsynced      int64 `json:"unsynced"`
	}
)

type (
	CreateProductRequest struct {
		Name        string     `json:"name" binding:"required,min=3,max=100" example:"Refined Bronze Hat"`
		Description string     `json:"description,omitempty" example:"A very nice hat"`
		SKU         string     `json:"sku" binding:"required" example:"L1L-448"`
		Stock       int        `json:"stock" binding:"required,min=0" example:"100"`
		CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	}

	UpdateProductRequest struct {
		ID			uuid.UUID	`json:"-"`
		Name        string     `json:"name" example:"Refined Bronze Hat"`
		Description *string    `json:"description,omitempty" example:"A very nice hat"`
		SKU         string     `json:"sku" example:"L1L-448"`
		Stock       int        `json:"stock" example:"100"`
		CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	}

	UpdateStockRequest struct {
		Change int    `json:"change" binding:"required" example:"-1"`
		Source string `json:"source" binding:"required" example:"shopee"`
		Note   string `json:"note" example:"Order #12345"`
	}

	ProductResponse struct {
		ID               uuid.UUID                 `json:"id"`
		Name             string                    `json:"name" example:"Refined Bronze Hat"`
		Description      string                    `json:"description"`
		SKU              string                    `json:"sku" example:"L1L-448"`
		Stock            int                       `json:"stock" example:"100"`
		Category         *ProductCategoryResponse  `json:"category,omitempty"`
		Images           []ProductImageResponse    `json:"images,omitempty"`
		ExternalProducts []ExternalProductResponse `json:"external_products,omitempty"`
		CreatedAt        time.Time                 `json:"created_at"`
		UpdatedAt        time.Time                 `json:"updated_at"`
	}

	ProductListResponse struct {
		ID               uuid.UUID                 `json:"id"`
		Name             string                    `json:"name"`
		SKU              string                    `json:"sku"`
		Stock            int                       `json:"stock"`
		Category         *ProductCategoryResponse  `json:"category,omitempty"`
		Images           []ProductImageResponse    `json:"images,omitempty"`
		ExternalProducts []ExternalProductResponse `json:"external_products,omitempty"`
		Status           string                    `json:"status"` // "Mapped", "Unmapped", "Low Stock", "Out of Stock"
		CreatedAt        time.Time                 `json:"created_at"`
	}

	ProductPaginationResponse struct {
		response.PaginationResponse
		Data []ProductListResponse `json:"data"`
	}
	ProductPaginationRepositoryResponse struct {
		response.PaginationResponse
		Products []entity.Product
	}
)

// Vendor
type (
	VendorResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		PhoneNumber string    `json:"phone"`
		Address     string    `json:"address"`
		ImageURL    string    `json:"image_url"`
		Description string    `json:"description"`
	}
	CreateVendorRequest struct {
		Name        string `json:"name" binding:"required,min=3,max=100"`
		Email       string `json:"email,omitempty" binding:"omitempty,email"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Address     string `json:"address,omitempty"`
		ImageURL    string `json:"image_url,omitempty"`
		Description string `json:"description,omitempty"`
	}
	UpdateVendorRequest struct {
		ID          uuid.UUID `json:"-"`
		Name        string    `json:"name" binding:"required,min=3,max=100"`
		Email       *string   `json:"email,omitempty" binding:"omitempty,email"`
		PhoneNumber string    `json:"phone_number" binding:"required"`
		Address     *string   `json:"address,omitempty" binding:"omitempty"`
		ImageURL    *string   `json:"image_url,omitempty" binding:"omitempty"`
		Description *string   `json:"description,omitempty" binding:"omitempty"`
	}
	VendorPaginationResponse struct {
		response.PaginationResponse
		Data []VendorResponse `json:"data"`
	}
	VendorPaginationRepositoryResponse struct {
		response.PaginationResponse
		Vendors []entity.Vendor
	}
)
