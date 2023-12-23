package router

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pdstuber/fiber-test/internal/api/handlers"
)

const shutdownTimeout = 5 * time.Second

type Server struct {
	listenPort string
	errChan    chan error
	fiberApp   *fiber.App
}

func New(listenPort string) *Server {
	app := fiber.New()

	app.Get("/", handlers.Hello)

	return &Server{
		listenPort: listenPort,
		fiberApp:   app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.fiberApp.Listen(s.listenPort); err != nil {
			s.errChan <- err
		}
	}()

	select {
	case err := <-s.errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop() error {
	return s.fiberApp.ShutdownWithTimeout(shutdownTimeout)
}
