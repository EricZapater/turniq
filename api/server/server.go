package server

import (
	"api/config"
	"api/internal/auth"
	"api/internal/customers"
	"api/internal/jobs"
	"api/internal/operators"
	"api/internal/payments"
	"api/internal/scheduleentries"
	"api/internal/shifts"
	"api/internal/shopfloors"
	"api/internal/timeentries"
	"api/internal/users"
	"api/internal/workcenters"
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
	customerRepo := customers.NewRepository(s.db)
	operatorRepo := operators.NewRepository(s.db)
	jobRepo := jobs.NewRepository(s.db)
	paymentRepo := payments.NewRepository(s.db)
	shopfloorRepo := shopfloors.NewRepository(s.db)
	workcenterRepo := workcenters.NewRepository(s.db)
	shiftRepo := shifts.NewRepository(s.db)
	scheduleEntryRepo := scheduleentries.NewRepository(s.db)
	timeEntryRepo := timeentries.NewRepository(s.db)

	//Services
	customerService := customers.NewService(customerRepo)
	userService := users.NewService(userRepo, customerService)
	authService := auth.NewAuthService(userService,customerService, authMiddleware)
	operatorService := operators.NewService(operatorRepo, customerService)
	jobService := jobs.NewService(jobRepo, customerService)
	paymentService := payments.NewService(paymentRepo)
	shopfloorService := shopfloors.NewService(shopfloorRepo, customerService)
	workcenterService := workcenters.NewService(workcenterRepo, customerService)
	shiftService := shifts.NewService(shiftRepo)
	scheduleEntryService := scheduleentries.NewService(scheduleEntryRepo)
	timeEntryService := timeentries.NewService(timeEntryRepo)
	//Handlers
	userHandler := users.NewHandler(userService)
	customerHandler := customers.NewHandler(customerService)
	authHandler := auth.NewAuthHandler(authService, authMiddleware)
	operatorHandler := operators.NewHandler(operatorService)
	jobHandler := jobs.NewHandler(jobService)
	paymentHandler := payments.NewHandler(paymentService)
	shopfloorHandler := shopfloors.NewHandler(shopfloorService)
	workcenterHandler := workcenters.NewHandler(workcenterService)
	shiftHandler := shifts.NewHandler(shiftService)
	scheduleEntryHandler := scheduleentries.NewHandler(scheduleEntryService)
	timeEntryHandler := timeentries.NewHandler(timeEntryService)
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// Prometheus metrics endpoint
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Public routes
	public := s.router.Group("/auth")
	auth.RegisterRoutes(public, authHandler, authMiddleware)
	users.RegisterAdminRoutes(public, &userHandler)

	//protected routes
	protected := s.router.Group("/api")
	protected.Use(authMiddleware.MiddlewareFunc())
	protected.Use(middleware.ContextMiddleware()) // Inject context values
	users.RegisterRoutes(protected, &userHandler)	
	customers.RegisterRoutes(protected, &customerHandler)
	operators.RegisterRoutes(protected, &operatorHandler)
	jobs.RegisterRoutes(protected, &jobHandler)
	payments.RegisterRoutes(protected, &paymentHandler)
	shopfloors.RegisterRoutes(protected, &shopfloorHandler)
	workcenters.RegisterRoutes(protected, &workcenterHandler)
	shifts.RegisterRoutes(protected, &shiftHandler)
	scheduleentries.RegisterRoutes(protected, &scheduleEntryHandler)
	timeentries.RegisterRoutes(protected, &timeEntryHandler)
	return nil
	
}

func(s *Server)Run()error{
	return s.router.Run(":" + s.config.App.Port)
}