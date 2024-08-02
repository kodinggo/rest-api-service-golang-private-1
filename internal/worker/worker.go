package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

func NewWorkerServer(redisOpt asynq.RedisClientOpt) error {
	// init asynq server
	svr := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 1,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			// log.Println("task ID", task.ResultWriter().TaskID())
			log.Println("task payload", string(task.Payload()))
		}),
	})

	mux := asynq.NewServeMux()

	task := taskHandler{}
	mux.HandleFunc(TaskSendEmail, task.sendEmailHandler)
	mux.HandleFunc(TaskUploadFile, task.uplaodFileHandler)

	return svr.Run(mux)
}

type WorkerClient struct {
	client *asynq.Client
}

func NewWorkerClient(redisOpt asynq.RedisClientOpt) *WorkerClient {
	client := asynq.NewClient(redisOpt)
	return &WorkerClient{client: client}
}

func (wrk WorkerClient) SendEmail(payload SendEmail) (*asynq.TaskInfo, error) {
	sendEmailPayload, _ := json.Marshal(payload)
	task := asynq.NewTask(TaskSendEmail, sendEmailPayload)
	return wrk.client.Enqueue(task)
}
