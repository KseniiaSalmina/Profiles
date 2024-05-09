package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/config"
)

type Storage interface {
	ReturnSalt() string
	GetAuthData(username string) (string, bool, error)
	GetAllUsers(offset, limit int) *models.PageUsers
	AddUser(user models.UserRequest) (string, error)
	GetUserByID(id string) (*models.UserResponse, error)
	ChangeUser(user models.UserRequest) error
	DeleteUser(id string) error
}

type Server struct {
	httpServer *http.Server
	storage    Storage
}

func NewServer(cfg config.Server, storage Storage) *Server {
	s := &Server{storage: storage}

	router := bunrouter.New().Compat()
	router.GET("/user", s.getAllUsers)
	router.POST("/user", s.postUser)
	router.GET("/user/:id", s.getUser)
	router.PATCH("/user/:id", s.patchUser)
	router.DELETE("/user/:id", s.deleteUser)

	//swagHandler := httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))
	//router.GET("/swagger/*path", swagHandler)

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
	log.Println("server started")

	go func() {
		err := s.httpServer.ListenAndServe()
		log.Printf("server stopped: %s", err.Error())
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
