package main

import (
	"log"
	"os"

	"github.com/Amierza/simponi-backend/cmd"
	"github.com/Amierza/simponi-backend/config/database"
	_ "github.com/Amierza/simponi-backend/docs"
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/logger"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/routes"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Simponi Backend API
//	@version		1.0
//	@description	REST API for Simponi - multi-store inventory and omnichannel management system.
//	@description	This API supports store management, product management, external product integration, vendor management, and role-based access control (RBAC).

//	@termsOfService	https://simponi.app/terms/

//	@contact.name	Simponi Support
//	@contact.email	support@simponi.app

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Use JWT token with format: Bearer <your_token>

//	@tag.name			Auth
//	@tag.description	Authentication (signin, refresh token) and token management

//	@tag.name			Users
//	@tag.description	User management and profile operations

//	@tag.name			Roles
//	@tag.description	Role management and RBAC configuration

//	@tag.name			Permissions
//	@tag.description	Permission management for access control

//	@tag.name			Stores
//	@tag.description	Store management and store-related operations

//	@tag.name			Store Users
//	@tag.description	Manage users within a store (membership & access)

//	@tag.name			Products
//	@tag.description	Product and inventory management

//	@tag.name			External Products
//	@tag.description	External product integration (Shopee, Tokopedia, etc.)

//	@tag.name			Vendors
//	@tag.description	Vendor and supplier management

//	@tag.name			Orders
//	@tag.description	Order and transaction management

//	@tag.name			Logs
//	@tag.description	System activity logs

//	@tag.name			Inventory Logs
//	@tag.description	Inventory movement and stock change logs

//	@tag.name			Uploads
//	@tag.description	File upload and asset management

//	@tag.name			Impersonate
//	@tag.description	Admin impersonation for debugging and support

// @externalDocs.description	OpenAPI Specification
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// Zap logger
	zapLogger, err := logger.New()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer zapLogger.Sync() // flush buffer

	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	var (
		// jwt
		jwt = jwt.NewJWT()

		// External API
		// externalGateway = gateway.NewExternalGateway(os.Getenv("API_EXTERNAL"), zapLogger)

		// Transaction
		tx = repository.NewTransaction(db)

		// Upload
		uploadService = service.NewUploadService(zapLogger)
		uploadHandler = handler.NewUploadHandler(uploadService, zapLogger)

		// Role Permission
		rolePermissionRepo = repository.NewRolePermissionRepository(db)

		// Permission
		permissionRepo    = repository.NewPermissionRepository(db)
		permissionService = service.NewPermissionService(permissionRepo, zapLogger, jwt)
		permissionHandler = handler.NewPermissionHandler(permissionService, zapLogger)

		// Role
		roleRepo    = repository.NewRoleRepository(db)
		roleService = service.NewRoleService(tx, roleRepo, permissionRepo, rolePermissionRepo, zapLogger, jwt)
		roleHandler = handler.NewRoleHandler(roleService, zapLogger)

		// User
		userRepo    = repository.NewUserRepository(db)
		userService = service.NewUserService(userRepo, roleRepo, zapLogger, jwt)
		userHandler = handler.NewUserHandler(userService, zapLogger)

		// Impersonate
		impersonateService = service.NewImpersonateService(userRepo, permissionRepo, zapLogger, jwt)
		impersonateHandler = handler.NewImpersonateHandler(impersonateService, zapLogger)

		// Authentication
		authService = service.NewAuthService(userRepo, permissionRepo, zapLogger, jwt)
		authHandler = handler.NewAuthHandler(authService, zapLogger)

		// Store User
		storeUserRepo    = repository.NewStoreUserRepository(db)
		storeUserService = service.NewStoreUserService(storeUserRepo, zapLogger, jwt)
		storeUserHandler = handler.NewStoreUserHandler(storeUserService, zapLogger)

		// Platform
		platformRepo = repository.NewPlatformRepository(db)

		// Store Platform
		storePlatformRepo = repository.NewStorePlatformRepository(db)

		// Store
		storeRepo    = repository.NewStoreRepository(db)
		storeService = service.NewStoreService(tx, storeRepo, storeUserRepo, platformRepo, storePlatformRepo, zapLogger, jwt)
		storeHandler = handler.NewStoreHandler(storeService, zapLogger)

		// Product Categories
		productCategoriesRepo    = repository.NewProductCategoriesRepository(db)
		productCategoriesService = service.NewProductCategoriesService(productCategoriesRepo, zapLogger, jwt)
		productCategoriesHandler = handler.NewProductCategoriesHandler(productCategoriesService, zapLogger)

		// Product
		productRepo    = repository.NewProductRepository(db)
		productService = service.NewProductService(productRepo, storeRepo, zapLogger, jwt)
		productHandler = handler.NewProductHandler(productService, zapLogger)

		// External Product
		externalProductRepo    = repository.NewExternalProductRepository(db)
		externalProductService = service.NewExternalProductService(externalProductRepo, productRepo, storeRepo, platformRepo, storePlatformRepo, zapLogger, jwt)
		externalProductHandler = handler.NewExternalProductHandler(externalProductService, zapLogger)

		// Vendor
		vendorRepo    = repository.NewVendorRepository(db)
		vendorService = service.NewVendorService(vendorRepo, zapLogger, jwt)
		vendorHandler = handler.NewVendorHandler(vendorService, zapLogger)

		// Order
		orderRepo    = repository.NewOrderRepository(db)
		orderService = service.NewOrderService(orderRepo, zapLogger, jwt)
		orderHandler = handler.NewOrderHandler(orderService, zapLogger)

		// Logging
		logRepo    = repository.NewLogRepository(db)
		logService = service.NewLogService(logRepo, zapLogger, jwt)
		logHandler = handler.NewLogHandler(logService, zapLogger)

		// Inventory Logging
		inventoryLogRepo    = repository.NewInventoryLoggingRepository(db)
		inventoryLogService = service.NewInventoryLoggingService(inventoryLogRepo, zapLogger, jwt)
		inventoryLogHandler = handler.NewInventoryLoggingHandler(inventoryLogService, zapLogger)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.Upload(server, uploadHandler, jwt, rolePermissionRepo)
	routes.Permission(server, permissionHandler, jwt, rolePermissionRepo)
	routes.Role(server, roleHandler, jwt, rolePermissionRepo)
	routes.User(server, userHandler, jwt, rolePermissionRepo)
	routes.Impersonate(server, impersonateHandler, jwt, rolePermissionRepo)
	routes.Auth(server, authHandler)
	routes.StoreUser(server, storeUserHandler, jwt, rolePermissionRepo)
	routes.Store(server, storeHandler, jwt, rolePermissionRepo)
	routes.ProductCategories(server, productCategoriesHandler, jwt)
	routes.Product(server, productHandler, jwt, rolePermissionRepo)
	routes.ExternalProduct(server, externalProductHandler, jwt, rolePermissionRepo)
	routes.Order(server, orderHandler, jwt, rolePermissionRepo)
	routes.Vendor(server, vendorHandler, jwt, rolePermissionRepo)
	routes.Log(server, logHandler, jwt, rolePermissionRepo)
	routes.InventoryLog(server, inventoryLogHandler, jwt, rolePermissionRepo)
	// swagger endpoint
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
