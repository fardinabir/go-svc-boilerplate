// Package server provides the API server for the application.
package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Server is the interface for the server
type Server interface {
	Name() string
	Run() error
	Shutdown(ctx context.Context) error
}

// apiServer is the server for backend API endpoints
type apiServer struct {
	port   int
	engine *echo.Echo
	log    *log.Entry
	db     *gorm.DB
}

func (s *apiServer) Name() string {
	return "apiServer"
}

// Run starts the User API server
func (s *apiServer) Run() error {
	log.Infof("%s serving on port %d", s.Name(), s.port)
	return s.engine.Start(fmt.Sprintf(":%d", s.port))
}

// Shutdown stops the User API server
func (s *apiServer) Shutdown(ctx context.Context) error {
	log.Infof("shutting down %s serving on port %d", s.Name(), s.port)
	return s.engine.Shutdown(ctx)
}
