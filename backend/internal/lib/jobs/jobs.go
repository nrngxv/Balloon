package job

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sriniously/go-boilerplate/internal/config"
)

type JobService struct {
	Client		*asynq.Client 
	Server		*asynq.Server
	logger		*zerolog.Logger
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	redisAddr := cfg.Redis.Address

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config {
			Concurrency: 10,
			Queues: map[string]int {
				"critical": 6,
				"default": 3,
				"low": 1,

			}
		}
	)

	return &JobService{
		Client: client,
		Server: server,
		logger: logger,
	}
}

func (j *JobService) Start() error {
	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)

	j.logger.Info().Msg("Starting background job server")
	if err := j.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping background job server")
	j.server.Shutdown()
	j.Client.Close()
}