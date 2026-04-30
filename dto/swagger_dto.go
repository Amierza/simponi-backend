package dto

// =====================
// COMMON - ERROR RESPONSE
// =====================
type ErrorResponse struct {
	Status  bool   `json:"status" example:"false"`
	Message string `json:"message" example:"failed request"`
	Error   string `json:"error" example:"invalid request"`
}

// =====================
// AUTHENTICATION
// =====================
type (
	SignInResponseWrapper struct {
		Status  bool           `json:"status" example:"true"`
		Message string         `json:"message" example:"success signin"`
		Data    SignInResponse `json:"data"`
	}
	RefreshTokenResponseWrapper struct {
		Status  bool                 `json:"status" example:"true"`
		Message string               `json:"message" example:"success refresh token"`
		Data    RefreshTokenResponse `json:"data"`
	}
)

// =====================
// USER - SUCCESS RESPONSE
// =====================

type (
	UserResponseWrapper struct {
		Status  bool         `json:"status" example:"true"`
		Message string       `json:"message" example:"success get user"`
		Data    UserResponse `json:"data"`
	}
	UsersResponseWrapper struct {
		Status  bool            `json:"status" example:"true"`
		Message string          `json:"message" example:"success get users"`
		Data    []*UserResponse `json:"data"`
		Meta    any             `json:"meta"`
	}
	UserEmptyResponseWrapper struct {
		Status  bool   `json:"status" example:"true"`
		Message string `json:"message" example:"success delete user"`
	}
)

// =====================
// ROLE - SUCCESS RESPONSE
// =====================

type (
	RoleResponseWrapper struct {
		Status  bool         `json:"status" example:"true"`
		Message string       `json:"message" example:"success get role"`
		Data    RoleResponse `json:"data"`
	}
	RolesResponseWrapper struct {
		Status  bool            `json:"status" example:"true"`
		Message string          `json:"message" example:"success get roles"`
		Data    []*RoleResponse `json:"data"`
		Meta    any             `json:"meta"`
	}
	RoleEmptyResponseWrapper struct {
		Status  bool   `json:"status" example:"true"`
		Message string `json:"message" example:"success delete role"`
	}
)

// =====================
// PERMISSION - SUCCESS RESPONSE
// =====================

type (
	PermissionResponseWrapper struct {
		Status  bool                  `json:"status" example:"true"`
		Message string                `json:"message" example:"success get permissions"`
		Data    []*PermissionResponse `json:"data"`
	}
	PermissionPaginationResponseWrapper struct {
		Status  bool                  `json:"status" example:"true"`
		Message string                `json:"message" example:"success get permissions"`
		Data    []*PermissionResponse `json:"data"`
		Meta    any                   `json:"meta"`
	}
)

// =====================
// STORE - SUCCESS RESPONSE
// =====================

type (
	StoreResponseWrapper struct {
		Status  bool          `json:"status" example:"true"`
		Message string        `json:"message" example:"success get store"`
		Data    StoreResponse `json:"data"`
	}
	StoresResponseWrapper struct {
		Status  bool             `json:"status" example:"true"`
		Message string           `json:"message" example:"success get stores"`
		Data    []*StoreResponse `json:"data"`
		Meta    any              `json:"meta"`
	}
	StoreEmptyResponseWrapper struct {
		Status  bool   `json:"status" example:"true"`
		Message string `json:"message" example:"success delete store"`
	}
)

// =====================
// STORE USER - SUCCESS RESPONSE
// =====================

type (
	StoreUsersResponseWrapper struct {
		Status  bool                  `json:"status" example:"true"`
		Message string                `json:"message" example:"success get store users"`
		Data    []*CustomUserResponse `json:"data"`
		Meta    any                   `json:"meta"`
	}
	StoreUserResponseWrapper struct {
		Status  bool               `json:"status" example:"true"`
		Message string             `json:"message" example:"success get store user"`
		Data    CustomUserResponse `json:"data"`
	}
	StoreUserEmptyResponseWrapper struct {
		Status  bool   `json:"status" example:"true"`
		Message string `json:"message" example:"success delete store user"`
	}
)

// =====================
// PRODUCT - SUCCESS RESPONSE
// =====================

type (
	ProductResponseWrapper struct {
		Status  bool            `json:"status" example:"true"`
		Message string          `json:"message" example:"success get product"`
		Data    ProductResponse `json:"data"`
	}
	ProductsResponseWrapper struct {
		Status  bool                  `json:"status" example:"true"`
		Message string                `json:"message" example:"success get products"`
		Data    []ProductListResponse `json:"data"`
		Meta    any                   `json:"meta"`
	}
	ProductStatsResponseWrapper struct {
		Status  bool                 `json:"status" example:"true"`
		Message string               `json:"message" example:"success get product stats"`
		Data    ProductStatsResponse `json:"data"`
	}
	ProductEmptyResponseWrapper struct {
		Status  bool   `json:"status" example:"true"`
		Message string `json:"message" example:"success"`
	}
)

// =======================
// External Product - SUCCESS RESPONSE
// =======================

type (
	ExternalProductSuccessResponse struct {
		Status  bool                    `json:"status" example:"true"`
		Message string                  `json:"message" example:"success get external product"`
		Data    ExternalProductResponse `json:"data"`
	}
	ExternalProductsSuccessResponse struct {
		Status  bool                      `json:"status" example:"true"`
		Message string                    `json:"message" example:"success get external products"`
		Data    []ExternalProductResponse `json:"data"`
	}
	ExternalProductCreateSuccessResponse struct {
		Status  bool                    `json:"status" example:"true"`
		Message string                  `json:"message" example:"success create external product"`
		Data    ExternalProductResponse `json:"data"`
	}
	ExternalProductUpdateSuccessResponse struct {
		Status  bool                    `json:"status" example:"true"`
		Message string                  `json:"message" example:"success update external product"`
		Data    ExternalProductResponse `json:"data"`
	}
	ExternalProductDeleteSuccessResponse struct {
		Status  bool   `json:"status" example:"true"`
		Message string `json:"message" example:"success delete external product"`
	}
)

// =======================
// Upload - SUCCESS RESPONSE
// =======================

type (
	UploadSuccessSingleResponse struct {
		Status  bool                `json:"status" example:"true"`
		Message string              `json:"message" example:"success upload file"`
		Data    UploadImageResponse `json:"data"`
	}
	UploadSuccessMultipleResponse struct {
		Status  bool                  `json:"status" example:"true"`
		Message string                `json:"message" example:"success upload files"`
		Data    []UploadImageResponse `json:"data"`
	}
)

// =======================
// Impersonate - SUCCESS RESPONSE
// =======================

type ImpersonateSuccessResponse struct {
	Status  bool                `json:"status" example:"true"`
	Message string              `json:"message" example:"success impersonate"`
	Data    ImpersonateResponse `json:"data"`
}
