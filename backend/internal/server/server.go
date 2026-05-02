package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nrngxv/go-boilerplate/internal/config"
	"github.com/nrngxv/go-boilerplate/internal/database"
	job "github.com/nrngxv/go-boilerplate/internal/lib/jobs"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type server struct {
	Config 			*config.Config
	Logger 			*zerolog.Logger
	LoggerService	*loggerPkg.LoggerService
	DB 				*database.Database
	Redis 			*redis.Client
	httpServer		*http.Server
	Job				*job.JobService
}

func New(cfg config.Config, logger *zerolog.Logger, loggerService *loggerPkg.loggerservice) (*Server, error) {
	db, err := database.New(cfg, logger, loggerService)
	if err != nil {
		return nil, fmt.Errorf("Database initialisation failed: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	if loggerService != nil && loggerService.GetApplication() != nil {
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	}

	//testing the redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobService := job.NewJobService(logger, cfg)
	jobService.InitHandlers(cfg, logger)

	if err := jobService.Start(); err != nil {
		return nil, err
	}

	server := &Server{
		Config:        cfg,
		Logger:        logger,
		LoggerService: loggerService,
		DB:            db,
		Redis:         redisClient,
		Job:           jobService,
	}

	return server, nil

}

//getting all the functionality for the server with whatever defined in the configs
func (s *server) setupHTTPServer(handler http.handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("env", s.Config.Primary.Env).
		Msg("starting server")

	return s.httpServer.ListenAndServe()
}


//method to shutdown the server. Gracefully.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	if s.Job != nil {
		s.Job.Stop()
	}

	return nil
}