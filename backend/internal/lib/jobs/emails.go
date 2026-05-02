package job

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynqf"
)

const (
	TaskWelcome = "email: welcome"
)

type WelcomeEmailTask struct {
	To			string
	Firstname	string
}

func NewWelcomeEmailTask(to, firstName string) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{
		To: 		to,
		FirstName: 	firstName,
	})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskWelcome, payload,
	asynq.MaxRetry(3),
	asynq.Queue("default"),
	asynq.Timeout(30*time.Second)), nil
	
}