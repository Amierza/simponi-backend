package main

import (
	"log"
	"os"

	"github.com/Amierza/simponi-backend/cmd"
	"github.com/Amierza/simponi-backend/config/database"
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/logger"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/routes"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
)

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

		// Resource
		// Upload
		// uploadService = service.NewUploadService(zapLogger)
		// uploadHandler = handler.NewUploadHandler(uploadService, zapLogger)

		// User
		userRepo    = repository.NewUserRepository(db)
		userService = service.NewUserService(userRepo, jwt, zapLogger)
		userHandler = handler.NewUserHandler(userService, zapLogger)

		// Authentication
		authService = service.NewAuthService(userRepo, zapLogger, jwt)
		authHandler = handler.NewAuthHandler(authService, zapLogger)

		// Permission
		permissionRepo    = repository.NewPermissionRepository(db)
		permissionService = service.NewPermissionService(permissionRepo)

		// Logging
		logRepo    = repository.NewLogRepository(db)
		logService = service.NewLogService(logRepo, zapLogger, jwt)
		logHandler = handler.NewLogHandler(logService, zapLogger)

		// Product
		productRepo    = repository.NewProductRepository(db)
		productService = service.NewProductService(productRepo, zapLogger, jwt)
		productHandler = handler.NewProductHandler(productService, zapLogger)

		// Vendor
		vendorRepo    = repository.NewVendorRepository(db)
		vendorService = service.NewVendorService(vendorRepo, zapLogger, jwt)
		vendorHandler = handler.NewVendorHandler(vendorService, zapLogger)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// routes.Upload(server, uploadHandler, jwt)
	routes.Auth(server, authHandler)
	routes.User(server, userHandler, jwt, permissionService)
	routes.Log(server, logHandler, jwt, permissionService)
	routes.Product(server, productHandler, jwt, permissionService)
	routes.Vendor(server, vendorHandler, jwt, permissionService)

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
