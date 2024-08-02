package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// define event type
const (
	TaskSendEmail  = "SEND_EMAIL"
	TaskUploadFile = "UPLOAD_FILE"
)

// define event payload
type (
	SendEmail struct {
		From    string
		To      string
		Subject string
		Content string
	}

	UplaodFile struct {
		Origin      string
		Destination string
	}
)

type taskHandler struct {
}

func (th taskHandler) sendEmailHandler(ctx context.Context, t *asynq.Task) error {
	time.Sleep(1 * time.Second)

	// Mendapatkan paylaod
	var payload SendEmail
	_ = json.Unmarshal(t.Payload(), &payload)

	log.Printf("Sending email, subject %s....", payload.Subject)
	return nil
}

func (th taskHandler) uplaodFileHandler(ctx context.Context, t *asynq.Task) error {
	log.Println("Uploading file....")
	return nil
}
