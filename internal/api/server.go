package api

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/uptrace/bunrouter"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/config"
	"github.com/KseniiaSalmina/Profiles/internal/database"
)

type Storage interface {
	ReturnSalt() string
	GetAuthData(username string) (*database.User, error)
	GetAllUsers(offset, limit int) *models.PageUsers
	AddUser(user models.UserRequest) (string, error)
	GetUserByID(id string) (*models.UserResponse, error)
	ChangeUser(user models.UserRequest) error
	DeleteUser(id string) error
}

type Server struct {
	httpServer *http.Server
	storage    Storage
	logger     *logrus.Logger
}

func NewServer(cfg config.Server, storage Storage, logger *logrus.Logger) *Server {
	s := &Server{storage: storage, logger: logger}

	router := bunrouter.New().Compat()
	router.GET("/user", s.getAllUsers)
	router.POST("/user", s.postUser)
	router.GET("/user/:id", s.getUser)
	router.PUT("/user/:id", s.putUser)
	router.DELETE("/user/:id", s.deleteUser)

	swagHandler := httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))
	router.GET("/swagger/*path", swagHandler)

	s.httpServer = &http.Server{
		Addr:         cfg.Listen,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return s
}

func (s *Server) Run() {
	s.logger.Info("server started")

	go func() {
		err := s.httpServer.ListenAndServe()
		s.logger.Infof("server stopped: %s", err.Error())
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
