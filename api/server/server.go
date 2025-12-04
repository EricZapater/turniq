package server

import (
	"api/config"
	"api/internal/auth"
	"api/internal/customers"
	"api/internal/jobs"
	"api/internal/operators"
	"api/internal/payments"
	"api/internal/tenants"
	"api/internal/users"
	"api/middleware"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router *gin.Engine
	config config.Config
	db *sql.DB
}

func NewServer(config config.Config, db *sql.DB) *Server {
	return &Server{
		router: gin.New(),
		config: config,
		db:     db,
	}
}

func(s *Server)Setup()error{
	s.router.Use(middleware.SetupCORS())
	s.router.Use(middleware.ObservabilityMiddleware())
	
	authMiddleware, err := middleware.SetupJWT(s.config)
	if err != nil{
		return err
	}
	

	//Repositories
	userRepo := users.NewRepository(s.db)
	tenantRepo := tenants.NewRepository(s.db)
	customerRepo := customers.NewRepository(s.db)
	operatorRepo := operators.NewRepository(s.db)
	jobRepo := jobs.NewRepository(s.db)
	paymentRepo := payments.NewRepository(s.db)

	//Services
	userService := users.NewService(userRepo)
	tenantService := tenants.NewService(tenantRepo)
	customerService := customers.NewService(customerRepo)
	authService := auth.NewAuthService(userRepo, authMiddleware)
	operatorService := operators.NewService(operatorRepo)
	jobService := jobs.NewService(jobRepo)
	paymentService := payments.NewService(paymentRepo)

	//Handlers
	userHandler := users.NewHandler(userService)
	tenantHandler := tenants.NewHandler(tenantService)
	customerHandler := customers.NewHandler(customerService)
	authHandler := auth.NewAuthHandler(authService, authMiddleware)
	operatorHandler := operators.NewHandler(operatorService)
	jobHandler := jobs.NewHandler(jobService)
	paymentHandler := payments.NewHandler(paymentService)
	

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// Prometheus metrics endpoint
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Public routes
	public := s.router.Group("/auth")
	auth.RegisterRoutes(public, authHandler, authMiddleware)

	//protected routes
	protected := s.router.Group("/api")
	protected.Use(authMiddleware.MiddlewareFunc())
	users.RegisterRoutes(protected, &userHandler)
	tenants.RegisterRoutes(protected, &tenantHandler)	
	customers.RegisterRoutes(protected, &customerHandler)
	operators.RegisterRoutes(protected, &operatorHandler)
	jobs.RegisterRoutes(protected, &jobHandler)
	payments.RegisterRoutes(protected, &paymentHandler)

	return nil
	
}

func(s *Server)Run()error{
	return s.router.Run(":" + s.config.App.Port)
}