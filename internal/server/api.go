// Package server provides the API server for the application.
package server

import (
	"fmt"

	"github.com/fardinabir/go-svc-boilerplate/internal/cases"
	"github.com/fardinabir/go-svc-boilerplate/internal/config"
	"github.com/fardinabir/go-svc-boilerplate/internal/db"
	"github.com/fardinabir/go-svc-boilerplate/internal/health"
	"github.com/fardinabir/go-svc-boilerplate/internal/user"
	"github.com/fardinabir/go-svc-boilerplate/pkg/logger"
	"github.com/fardinabir/go-svc-boilerplate/pkg/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

// APIServerOpts is the options for the API server
type APIServerOpts struct {
	ListenPort int
	Config     config.Config
}

// Handlers holds the per-domain HTTP handlers assembled by Wire.
type Handlers struct {
	User  user.Handler
	Cases cases.Handler
}

// NewAPI returns a new instance of the API server
func NewAPI(opts APIServerOpts) (Server, error) {
	logEntry := log.NewEntry(log.StandardLogger())
	log.SetFormatter(&log.JSONFormatter{})

	// Initialize global logger
	logger.InitLogger(logEntry)

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

	s := &apiServer{
		port:   opts.ListenPort,
		engine: engine,
		log:    logEntry,
		db:     dbInstance,
	}

	if err := s.setupRoutes(engine); err != nil {
		return nil, err
	}

	engine.Use(requestLogger())

	return s, nil
}

// setupRoutes builds the dependency graph via Wire and registers each domain's routes.
// Adding a domain is two lines here: add its ProviderSet to InitializeHandlers (wire.go)
// and a RegisterRoutes call below.
func (s *apiServer) setupRoutes(e *echo.Echo) error {
	// Each domain that registers custom validation tags must be listed here.
	// When adding a new domain with custom tags, append its RegisterValidations fn.
	e.Validator = web.NewCustomValidator(
		user.RegisterValidations,
		// cases.RegisterValidations,  // add when cases defines custom tags
	)

	handlers, err := InitializeHandlers(s.db)
	if err != nil {
		return fmt.Errorf("failed to initialize handlers: %w", err)
	}

	api := e.Group("/api/v1")

	health.RegisterRoutes(api, health.New())
	user.RegisterRoutes(api.Group("/users"), handlers.User)
	cases.RegisterRoutes(api.Group("/cases"), handlers.Cases)

	return nil
}
