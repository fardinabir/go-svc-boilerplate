package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// userAPIServer is the API server for User
type userAPIServer struct {
	port   int
	engine *echo.Echo
	log    *log.Entry
	db     *gorm.DB
}

func (s *userAPIServer) Name() string {
	return "userAPIServer"
}

// Run starts the User API server
func (s *userAPIServer) Run() error {
	log.Infof("%s serving on port %d", s.Name(), s.port)
	return s.engine.Start(fmt.Sprintf(":%d", s.port))
}

// Shutdown stops the User API server
func (s *userAPIServer) Shutdown(ctx context.Context) error {
	log.Infof("shutting down %s serving on port %d", s.Name(), s.port)
	return s.engine.Shutdown(ctx)
}
