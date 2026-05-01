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

	// Inventory Logging Errors
	FAILED_CREATE_INVENTORY_LOG          = "failed to create inventory log"
	FAILED_GET_INVENTORY_LOGS            = "failed to get inventory logs"
	FAILED_GET_INVENTORY_LOGS_BY_PRODUCT = "failed to get inventory logs by product ID"

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

	// Inventory Logging Success
	SUCCESS_CREATE_INVENTORY_LOG          = "success to create inventory log"
	SUCCESS_GET_INVENTORY_LOGS            = "success to get inventory logs"
	SUCCESS_GET_INVENTORY_LOGS_BY_PRODUCT = "success to get inventory logs by product ID"

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
	ErrForbidden     = errors.New("forbidden")

	// Input

	// Authentication
	ErrIncorrectPassword = errors.New("credential incorrect")

	// Logging
	ErrCreateLog          = errors.New("failed to create log")
	ErrGetLogs            = errors.New("failed to get logs")
	ErrGetLogsByStoreID   = errors.New("failed to get logs by store ID")
	ErrGetLogsByDateRange = errors.New("failed to get logs by date range")

	// Inventory Logging
	ErrCreateInventoryLog          = errors.New("failed to create inventory log")
	ErrGetInventoryLogs            = errors.New("failed to get inventory logs")
	ErrGetInventoryLogsByProductID = errors.New("failed to get inventory logs by product ID")

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

// Without Pagination
type (
	UploadImageResponse struct {
		ImageID  uuid.UUID `json:"image_id"`
		ImageURL string    `json:"image_url"`
	}

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

	// Impersonate
	ImpersonateResponse struct {
		AccessToken string `json:"access_token" example:"<new_access_token_here>"`
	}
)

// Role
type (
	RoleResponse struct {
		ID          uuid.UUID            `json:"id"`
		Name        string               `json:"name"`
		Permissions []PermissionResponse `json:"permissions"`
	}
	CreateRoleRequest struct {
		Name           string       `json:"name" binding:"required"`
		PermissionsIDs []*uuid.UUID `json:"permission_ids" binding:"required"`
	}
	UpdateRoleRequest struct {
		ID            uuid.UUID    `json:"-"`
		Name          string       `json:"name" binding:"required"`
		PermissionIDs []*uuid.UUID `json:"permission_ids" binding:"required"`
	}
	RolePaginationResponse struct {
		response.PaginationResponse
		Data []*RoleResponse `json:"data"`
	}
	RolePaginationRepositoryResponse struct {
		response.PaginationResponse
		Roles []*entity.Role
	}
)

// Permission
type (
	PermissionResponse struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Endpoint string    `json:"endpoint"`
		Method   string    `json:"method"`
		Module   string    `json:"module"`
	}
	PermissionPaginationResponse struct {
		response.PaginationResponse
		Data []*PermissionResponse `json:"data"`
	}
	PermissionPaginationRepositoryResponse struct {
		response.PaginationResponse
		Permissions []*entity.Permission
	}
)

// User
type (
	UserResponse struct {
		ID       uuid.UUID    `json:"id"`
		Email    string       `json:"email"`
		Name     string       `json:"name"`
		ImageURL string       `json:"image_url"`
		Status   string       `json:"status"`
		Role     RoleResponse `json:"role"`
	}
	CustomUserResponse struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}
	CreateUserRequest struct {
		Email    string     `json:"email" binding:"email"`
		Password string     `json:"password" binding:"required,min=5,max=8"`
		Name     string     `json:"name" binding:"required,min=3,max=100"`
		ImageURL string     `json:"image_url"`
		RoleID   *uuid.UUID `json:"role_id"`
	}
	UpdateUserRequest struct {
		ID       uuid.UUID  `json:"-"`
		Email    string     `json:"email" binding:"email"`
		Name     string     `json:"name" binding:"required,min=3,max=100"`
		ImageURL *string    `json:"image_url,omitempty" binding:"omitempty"`
		RoleID   *uuid.UUID `json:"role_id"`
	}
	UpdateUserStatus struct {
		ID     uuid.UUID `json:"-"`
		Status string    `json:"status" binding:"required"`
	}
	UserPaginationResponse struct {
		response.PaginationResponse
		Data []*UserResponse `json:"data"`
	}
	UserPaginationRepositoryResponse struct {
		response.PaginationResponse
		Users []*entity.User
	}
)

// Store User
type (
	CreateStoreUsersRequest struct {
		StoreID *uuid.UUID   `json:"-"`
		UserIDs []*uuid.UUID `json:"user_ids" binding:"required"`
	}
	StoreUserPaginationResponse struct {
		response.PaginationResponse
		Data []*CustomUserResponse `json:"data"`
	}
	StoreUserPaginationRepositoryResponse struct {
		response.PaginationResponse
		StoreUsers []*entity.StoreUser
	}
)

// Platform
type (
	PlatformResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
)

// Store
type (
	StoreResponse struct {
		ID          uuid.UUID          `json:"id"`
		Name        string             `json:"name"`
		Description string             `json:"description"`
		ImageURL    string             `json:"image_url"`
		IsActive    bool               `json:"is_active"`
		Platforms   []PlatformResponse `json:"platforms"`
	}
	CustomStoreResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	CreateStoreRequest struct {
		UserID      *uuid.UUID `json:"-"`
		PlatformID  *uuid.UUID `json:"platform_id" binding:"required"`
		Name        string     `json:"name" binding:"required,min=3,max=100"`
		ImageURL    string     `json:"image_url,omitempty"`
		Description string     `json:"description,omitempty"`
	}
	UpdateStoreRequest struct {
		ID          uuid.UUID `json:"-"`
		Name        string    `json:"name" binding:"required,min=3,max=100"`
		ImageURL    *string   `json:"image_url,omitempty" binding:"omitempty"`
		Description *string   `json:"description,omitempty" binding:"omitempty"`
		IsActive    *bool     `json:"is_active,omitempty" binding:"omitempty"`
	}
	StorePaginationResponse struct {
		response.PaginationResponse
		Data []*StoreResponse `json:"data"`
	}
	StorePaginationRepositoryResponse struct {
		response.PaginationResponse
		Stores []*entity.Store
	}
)

// Product
type (
	ProductResponse struct {
		ID               uuid.UUID                 `json:"id"`
		Name             string                    `json:"name" example:"Refined Bronze Hat"`
		Description      string                    `json:"description"`
		SKU              string                    `json:"sku" example:"L1L-448"`
		Stock            int                       `json:"stock" example:"100"`
		Store            *ProductStoreResponse     `json:"store,omitempty"`
		Category         *ProductCategoryResponse  `json:"category,omitempty"`
		Images           []ProductImageResponse    `json:"images,omitempty"`
		ExternalProducts []ExternalProductResponse `json:"external_products,omitempty"`
		CreatedAt        time.Time                 `json:"created_at"`
		UpdatedAt        time.Time                 `json:"updated_at"`
	}
	CreateProductRequest struct {
		StoreID     *uuid.UUID `json:"-"`
		Images     []string `json:"images" binding:"omitempty,dive"`
		CategoryID  *uuid.UUID `json:"category_id,omitempty"`
		Name        string     `json:"name" binding:"required,min=3,max=100" example:"Refined Bronze Hat"`
		Description string     `json:"description,omitempty" example:"A very nice hat"`
		SKU         string     `json:"sku" binding:"required" example:"L1L-448"`
		Stock       int        `json:"stock" binding:"required,min=0" example:"100"`
	}
	UpdateProductRequest struct {
		ID          uuid.UUID  `json:"-"`
		Name        string     `json:"name" example:"Refined Bronze Hat"`
		Description *string    `json:"description,omitempty" example:"A very nice hat"`
		SKU         string     `json:"sku" example:"L1L-448"`
		Stock       int        `json:"stock" example:"100"`
		CategoryID  *uuid.UUID `json:"category_id,omitempty"`
		StoreID     *uuid.UUID `json:"-,omitempty"`
	}
	UpdateStockRequest struct {
		ID      uuid.UUID  `json:"-"`
		Change  int        `json:"change" binding:"required" example:"-1"`
		Source  string     `json:"source" binding:"required" example:"shopee"`
		Note    string     `json:"note" example:"Order #12345"`
		StoreID *uuid.UUID `json:"-,omitempty"`
	}
	ProductListResponse struct {
		ID               uuid.UUID                 `json:"id"`
		Name             string                    `json:"name"`
		SKU              string                    `json:"sku"`
		Stock            int                       `json:"stock"`
		Store            *ProductStoreResponse     `json:"store,omitempty"`
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

	// Product Store
	ProductStoreResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

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

	// Product Image
	ProductImageResponse struct {
		ID       uuid.UUID `json:"id"`
		ImageURL string    `json:"image_url" example:"https://example.com/image.jpg"`
	}

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

// External Product
type (
	ExternalProductResponse struct {
		ID                uuid.UUID `json:"id"`
		ImageURL          string    `json:"image"`
		ProductName       string    `json:"product_name"`
		Platform          string    `json:"platform"`
		StorePlatformName string    `json:"store_platform_name"`
		Price             int64     `json:"price"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	}
	CreateExternalProductRequest struct {
		StoreID    *uuid.UUID `json:"-"`
		ProductID  *uuid.UUID `json:"product_id" binding:"required"`
		PlatformID *uuid.UUID `json:"platform_id" binding:"required"`
		Price      int64      `json:"price" binding:"required,min=0"`
	}
	UpdateExternalProductRequest struct {
		ID      uuid.UUID  `json:"-"`
		StoreID *uuid.UUID `json:"-"`
		Price   int64      `json:"price" binding:"required,min=0"`
	}
	ExternalProductPaginationResponse struct {
		response.PaginationResponse
		Data []ExternalProductResponse `json:"data"`
	}
	ExternalProductPaginationRepositoryResponse struct {
		response.PaginationResponse
		ExternalProducts []entity.ExternalProduct
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
		Data []*VendorResponse `json:"data"`
	}
	VendorPaginationRepositoryResponse struct {
		response.PaginationResponse
		Vendors []*entity.Vendor
	}
)

type (
	// Order
	OrderResponse struct {
		ID               uuid.UUID             `json:"id"`
		ExternalOrderID  string                `json:"external_order_id" example:"1234567890"`
		Ordernumber      string                `json:"order_number" example:"ORD-20230901-001"`
		StoreID          *uuid.UUID            `json:"store_id,omitempty"`
		Platform         string                `json:"platform,omitempty"`
		BuyerName        string                `json:"buyer_name" example:"John Doe"`
		BuyerEmail       string                `json:"buyer_email" example:"rizkyardiansyah@gmail.com"`
		BuyerPhone       string                `json:"buyer_phone" example:"+1234567890"`
		ReceipentName    string                `json:"receipent_name" example:"Jane Doe"`
		ReceipentPhone   string                `json:"receipent_phone" example:"+0987654321"`
		ShippingAddress  string                `json:"shipping_address" example:"123 Main St, Anytown, USA"`
		ShippingCity     string                `json:"shipping_city" example:"Anytown"`
		ShippingProvince string                `json:"shipping_province" example:"Anystate"`
		ShippingPostal   string                `json:"shipping_postal" example:"12345"`
		ShippingMethod   string                `json:"shipping_method" example:"JNE Regular"`
		TrackingNumber   string                `json:"tracking_number" example:"JNE1234567890"`
		SubtotalAmount   int64                 `json:"subtotal_amount" example:"100000"`
		ShippingFee      int64                 `json:"shipping_fee" example:"15000"`
		MarketplaceFee   int64                 `json:"marketplace_fee" example:"5000"`
		DiscountAmount   int64                 `json:"discount_amount" example:"10000"`
		TaxAmount        int64                 `json:"tax_amount" example:"10000"`
		TotalAmount      int64                 `json:"total_amount" example:"105000"`
		NetAmount        int64                 `json:"net_amount" example:"90000"`
		OrderStatus      string                `json:"order_status" example:"PENDING"`
		PaymentStatus    string                `json:"payment_status" example:"UNPAID"`
		PaymentMethod    string                `json:"payment_method" example:"Credit Card"`
		OrderedAt        *time.Time            `json:"ordered_at,omitempty" example:"2023-09-01T12:00:00Z"`
		PaidAt           *time.Time            `json:"paid_at,omitempty" example:"2023-09-01T12:30:00Z"`
		ShippedAt        *time.Time            `json:"shipped_at,omitempty" example:"2023-09-02T08:00:00Z"`
		CompletedAt      *time.Time            `json:"completed_at,omitempty" example:"2023-09-05T17:00:00Z"`
		CancelledAt      *time.Time            `json:"cancelled_at,omitempty" example:"2023-09-03T10:00:00Z"`
		OrderDetails     []OrderDetailResponse `json:"order_details,omitempty"`
		CreatedAt        time.Time             `json:"created_at"`
	}

	OrderDetailResponse struct {
		ID                uuid.UUID  `json:"id"`
		OrderID           *uuid.UUID `json:"order_id,omitempty"`
		ExternalProductID *uuid.UUID `json:"external_product_id,omitempty"`
		Quantity          int        `json:"quantity" example:"2"`

		Order           *OrderResponse           `json:"order,omitempty"`
		ExternalProduct *ExternalProductResponse `json:"external_product,omitempty"`
	}

	OrderPaginationResponse struct {
		response.PaginationResponse
		Data []OrderResponse `json:"data"`
	}

	OrderPaginationRepositoryResponse struct {
		Orders []entity.Order `json:"orders"`
		response.PaginationResponse
	}
)

// Logging
type (
	LogRequest struct {
		StoreID *uuid.UUID `json:"store_id" binding:"required"`
		Action  string     `json:"action" binding:"required" example:"Create"`
		Message string     `json:"message" binding:"required" example:"Created a new store"`
	}
	LogResponse struct {
		ID        uuid.UUID  `json:"id"`
		StoreID   *uuid.UUID `json:"store_id,omitempty"`
		Action    string     `json:"action" example:"Create"`
		Message   string     `json:"message" example:"Created a new store"`
		CreatedAt time.Time  `json:"created_at"`
	}
	LogPaginationResponse struct {
		response.PaginationResponse
		Data []LogResponse `json:"data"`
	}
	LogPaginationRepositoryResponse struct {
		response.PaginationResponse
		Logs []entity.Log
	}
)

// Inventory Log
type (
	InventoryLogRequest struct {
		ProductID *uuid.UUID `gorm:"type:uuid" json:"product_id"`
		Change    int        `json:"change"` // -1, -2, +10
		Source    string     `json:"source"` // shopee, tiktok, manual
		Note      string     `json:"note"`
	}
	InventoryLogResponse struct {
		ID        uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
		Product   *ProductResponse `gorm:"type:uuid" json:"product"`
		Change    int              `json:"change"` // -1, -2, +10
		Source    string           `json:"source"` // shopee, tiktok, manual
		Note      string           `json:"note"`
		CreatedAt time.Time        `json:"created_at"`
	}
	InventoryLogPaginationResponse struct {
		response.PaginationResponse
		Data []InventoryLogResponse `json:"data"`
	}
	InventoryLogPaginationRepositoryResponse struct {
		response.PaginationResponse
		InventoryLogs []entity.InventoryLog
	}
)
