package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"kasirpro/internal/config"
	"kasirpro/internal/database"
	"kasirpro/internal/handler"
	"kasirpro/internal/middleware"
	"kasirpro/internal/repository"
	"kasirpro/internal/service"
	"kasirpro/internal/utils"
)

func SetupApp() *fiber.App {
	cfg := config.LoadConfig()
	db, dbErr := database.InitDB()
	if db != nil && dbErr == nil {
		database.SeedData(db)
	}

	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ServerHeader: "KasirPro Engine",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return utils.ErrorResponse(c, code, err.Error())
		},
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(middleware.SecurityMiddleware())
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.StripServerHeader)

	app.Static("/uploads", "./uploads")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "online",
			"app":     cfg.AppName,
			"health":  "/api/v1/health",
			"message": "KasirPro Backend API is running",
		})
	})

	api := app.Group("/api/v1")

	// Health check endpoint (always responds even if DB connection failed)
	api.Get("/health", func(c *fiber.Ctx) error {
		if dbErr != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"message": "Database belum terhubung",
				"error":   dbErr.Error(),
				"db_host": cfg.DBHost,
				"app":     cfg.AppName,
			})
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "KasirPro backend is healthy and running", fiber.Map{
			"version":  "1.0.0",
			"app":      cfg.AppName,
			"database": "connected",
		})
	})

	// If database failed to connect, return 503 error for all functional endpoints
	if dbErr != nil {
		api.All("/*", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"message": "Database tidak dapat diakses",
				"error":   dbErr.Error(),
				"db_host": cfg.DBHost,
			})
		})
		return app
	}

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	transRepo := repository.NewTransactionRepository(db)

	authService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)
	customerService := service.NewCustomerService(customerRepo, userRepo)
	supplierService := service.NewSupplierService(supplierRepo)
	transService := service.NewTransactionService(transRepo, productRepo, customerRepo)
	midtransService := service.NewMidtransPaymentService(transRepo, customerRepo)

	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customerService)
	supplierHandler := handler.NewSupplierHandler(supplierService)
	transHandler := handler.NewTransactionHandler(transService, midtransService)
	dashboardHandler := handler.NewDashboardHandler(db)

	authGroup := api.Group("/auth", middleware.AuthRateLimiter())
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/register-staff", authHandler.RegisterStaff)
	authGroup.Post("/login", authHandler.Login)

	api.Get("/public/outlet-info/:code", customerHandler.GetPublicOutletInfo)
	api.Post("/public/member-register", customerHandler.RegisterPublicMember)

	protected := api.Group("", middleware.AuthMiddleware(), middleware.OutletIsolation(), middleware.GeneralAPIRateLimiter())

	protected.Get("/auth/me", authHandler.GetMe)
	protected.Post("/auth/logout", authHandler.Logout)

	protected.Get("/dashboard/summary", dashboardHandler.GetSummary)

	protected.Get("/categories", categoryHandler.GetAll)
	protected.Post("/categories", middleware.RequireRole("owner", "admin"), categoryHandler.Create)
	protected.Put("/categories/:id", middleware.RequireRole("owner", "admin"), categoryHandler.Update)
	protected.Delete("/categories/:id", middleware.RequireRole("owner", "admin"), categoryHandler.Delete)

	protected.Get("/products", productHandler.GetAll)
	protected.Get("/products/:id", productHandler.GetByID)
	protected.Get("/products/barcode/:barcode", productHandler.GetByBarcode)
	protected.Post("/products", middleware.RequireRole("owner", "admin"), productHandler.Create)
	protected.Put("/products/:id", middleware.RequireRole("owner", "admin"), productHandler.Update)
	protected.Delete("/products/:id", middleware.RequireRole("owner", "admin"), productHandler.Delete)

	protected.Get("/customers", customerHandler.GetAll)
	protected.Get("/customers/lookup/:code", customerHandler.LookupMember)
	protected.Post("/customers", customerHandler.Create)
	protected.Put("/customers/:id", middleware.RequireRole("owner", "admin"), customerHandler.Update)
	protected.Delete("/customers/:id", middleware.RequireRole("owner", "admin"), customerHandler.Delete)

	protected.Get("/suppliers", middleware.RequireRole("owner", "admin"), supplierHandler.GetAll)
	protected.Post("/suppliers", middleware.RequireRole("owner", "admin"), supplierHandler.Create)
	protected.Put("/suppliers/:id", middleware.RequireRole("owner", "admin"), supplierHandler.Update)
	protected.Delete("/suppliers/:id", middleware.RequireRole("owner", "admin"), supplierHandler.Delete)

	protected.Post("/pos/checkout", middleware.POSRateLimiter(), middleware.DiscountAuthorizationPolicy(), transHandler.Checkout)
	protected.Post("/pos/qris/confirm/:invoice", transHandler.ConfirmQRIS)
	protected.Get("/transactions", transHandler.GetAll)
	protected.Get("/transactions/:id", transHandler.GetByID)
	protected.Get("/transactions/invoice/:invoice", transHandler.GetByInvoice)
	protected.Post("/transactions/:id/void", middleware.RequireRole("owner", "admin"), transHandler.VoidTransaction)

	return app
}
