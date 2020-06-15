package worker

import (
	"fmt"
	"log"

	"github.com/squadcast_assignment/internal/infrastructure/database"
	queue "github.com/squadcast_assignment/internal/infrastructure/workqueue"
)

type WorkerService interface {
	ProcessWork(queueName string)
}

type worker struct {
	name string
	q    queue.QueueClient
}

func CreateWorker(name string, db database.DBClient, q queue.QueueClient) (WorkerService, error) {
	switch name {
	case "slack":
		log.Println("Creating slack worker")
		return slackWorker(name, q), nil
	case "jira":
		log.Println("Creating jira worker")
		return jiraWorker(name, db, q), nil
	case "zendesk":
		log.Println("Creating zendesk worker")
		return zendeskWorker(name, q), nil
	default:
		return nil, fmt.Errorf("worker: %v not found", name)
	}
}
