// Package server provides the API server for the application.
package server

import (
	"fmt"

	"github.com/fardinabir/go-svc-boilerplate/internal/controller"
	"github.com/fardinabir/go-svc-boilerplate/internal/db"
	"github.com/fardinabir/go-svc-boilerplate/internal/model"
	"github.com/fardinabir/go-svc-boilerplate/internal/repository"
	"github.com/fardinabir/go-svc-boilerplate/internal/service"
	"github.com/fardinabir/go-svc-boilerplate/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

// TxnAPIServerOpts is the options for the TxnAPIServer
type TxnAPIServerOpts struct {
	ListenPort int
	Config     model.Config
}

// NewAPI returns a new instance of the Txn API server
func NewAPI(opts TxnAPIServerOpts) (Server, error) {
	logger := log.NewEntry(log.StandardLogger())
	log.SetFormatter(&log.JSONFormatter{})

	// Initialize global logger
	utils.InitLogger(logger)

	dbInstance, err := db.New(opts.Config.PostgreSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	engine := echo.New()

	// Allow all origins for CORS
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	s := &userAPIServer{
		port:   opts.ListenPort,
		engine: engine,
		log:    logger,
		db:     dbInstance,
	}

	s.setupRoutes(engine)

	engine.Use(requestLogger())

	return s, nil
}

// initUserController creates and configures the user handler with its dependencies
//
//	Repository ====> Service =====> Controller
//
// It follows the CSR dependency injection pattern
func (s *userAPIServer) initUserController() controller.UserHandler {
	// Initialize dependencies (Repository -> Service -> Controller)
	userRepo := repository.NewUserRepository(s.db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserHandler(userService)

	return userController
}

// setupRoutes registers the routes for the application.
func (s *userAPIServer) setupRoutes(e *echo.Echo) {
	e.Validator = controller.NewCustomValidator()

	api := e.Group("/api/v1")

	// Health check
	healthHandler := controller.NewHealth()
	api.GET("/health", healthHandler.Health)

	userHandler := s.initUserController()

	controller.InitRoutes(api, userHandler)
}
