package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/squadcast_assignment/internal/infrastructure/database"
	queue "github.com/squadcast_assignment/internal/infrastructure/workqueue"
)

type jWorker struct {
	worker
	database.DBClient
}

func jiraWorker(name string, db database.DBClient, q queue.QueueClient) *jWorker {
	jWorker := jWorker{}
	jWorker.name = name
	jWorker.DBClient = db
	jWorker.q = q
	return &jWorker
}

func (jw *jWorker) ProcessWork(queueName string) {
	msgs, close, err := jw.q.Subscribe(queueName)
	if err != nil {
		panic(err)
	}
	defer close()

	stop := make(chan struct{})

	type incidentReport struct {
		ID             int64  `json:"id"`
		IncidentStatus string `json:"incident_status"`
	}

	go func() {
		for d := range msgs {
			iReport := &incidentReport{}
			if err := json.Unmarshal(d.Body, &iReport); err != nil {
				log.Println(err)
				continue
			}
			if iReport.IncidentStatus == "INCIDENT_CREATED" {
				timer := time.AfterFunc(time.Duration(5)*time.Minute, func() {
					row, err := jw.DBClient.Query(
						fmt.Sprintf(`SELECT message, ack, status FROM incidents WHERE id=%d;`, iReport.ID))
					if err != nil {
						log.Println(err)
						return
					}
					var message, ack, status string
					if row.Next() {
						row.Scan(&message, &ack, &status)
					}
					if message != "" && len(message) > 0 {
						if ack == "no" || status == "unresolved" {
							log.Printf(`A JIRA ticket with message: "%s" has been successfully created`, message)
						}
					}
				})
				defer timer.Stop()
			} else {
				log.Printf("%s with id: %d for JIRA", iReport.IncidentStatus, iReport.ID)
			}
			d.Ack(true)
		}
	}()
	log.Println("To exit press CTRL+C")
	<-stop
}
