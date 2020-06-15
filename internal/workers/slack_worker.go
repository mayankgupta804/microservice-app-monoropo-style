package worker

import (
	"encoding/json"
	"log"

	queue "github.com/squadcast_assignment/internal/infrastructure/workqueue"
)

type sWorker struct {
	worker
}

func slackWorker(name string, q queue.QueueClient) *sWorker {
	sWorker := sWorker{}
	sWorker.name = name
	sWorker.q = q
	return &sWorker
}

func (w *sWorker) ProcessWork(queueName string) {
	msgs, close, err := w.q.Subscribe(queueName)
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
			log.Printf("%s with id: %d for SLACK", iReport.IncidentStatus, iReport.ID)
			d.Ack(true)
		}
	}()
	log.Println("To exit press CTRL+C")
	<-stop
}
