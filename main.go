package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/squadcast_assignment/internal/config"
	client "github.com/squadcast_assignment/internal/eventhandler/grpc_client"
	server "github.com/squadcast_assignment/internal/eventhandler/grpc_server"
	"github.com/squadcast_assignment/internal/eventhandler/proto"
	"github.com/squadcast_assignment/internal/infrastructure/database"
	queue "github.com/squadcast_assignment/internal/infrastructure/workqueue"
	"github.com/squadcast_assignment/internal/migrations"
	"github.com/squadcast_assignment/internal/repository"
	"github.com/squadcast_assignment/internal/service"
	"github.com/squadcast_assignment/internal/webserver"
	worker "github.com/squadcast_assignment/internal/workers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.Load()

	clientApp := cli.NewApp()
	clientApp.Name = "Squadcast service"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "start:webserver",
			Description: "Start Incident Web Service",
			Action: func(c *cli.Context) error {
				return StartIncidentWebService()
			},
		},
		{
			Name:        "start:eventhandler",
			Description: "Start Event Handler Service",
			Action: func(c *cli.Context) error {
				return StartEventHandler()
			},
		},
		{
			Name:        "start:worker",
			Description: "Start worker of given type",
			Action: func(c *cli.Context) error {
				return StartWorker(c.Args().First())
			},
		},
		{
			Name:        "db:migrate:up",
			Description: "Create migrations",
			Action: func(c *cli.Context) error {
				return migrations.Up(config.App.Database)
			},
		},
		{
			Name:        "db:migrate:down",
			Description: "Destroy migrations",
			Action: func(c *cli.Context) error {
				return migrations.Down(config.App.Database)
			},
		},
	}
	if err := clientApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func StartIncidentWebService() error {
	var err error

	err = SetLogOutput(config.App.Logger, "incident_webservice.log")
	if err != nil {
		return fmt.Errorf("error encountered when setting log output: %v", err)
	}

	var db database.DBClient
	var repo repository.IncidentRepository
	var incidentService service.IncidentService
	var grpcClient proto.EventHandlerClient

	db, err = database.InitDatabaseConnection(config.App.Database)
	if err != nil {
		return fmt.Errorf("error encountered when connecting to DB: %v", err)
	}
	log.Println("successfully connected to database...")

	grpcClient, err = client.GetGRPClient(config.App.GRPCServer)
	if err != nil {
		return fmt.Errorf("error encountered when creating GRPC client: %v", err)
	}
	log.Println("grpc client created successfully...")

	defer func() {
		db.Close()
	}()

	repo = repository.InitIncidentRepository(db)
	incidentService = service.NewIncidentService(repo)
	router := webserver.SetupRoutes(incidentService, grpcClient)
	if err = webserver.StartServer(router, config.App.Server); err != nil {
		return err
	}
	return nil
}

func SetLogOutput(logConfig config.Logger, logFilename string) error {
	if logConfig.Output == "file" {
		f, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer f.Close()
		log.SetOutput(f)
	}
	return nil
}

func StartEventHandler() error {
	var err error
	var q queue.QueueClient

	err = SetLogOutput(config.App.Logger, "eventhandler.log")
	if err != nil {
		return fmt.Errorf("error encountered when setting log output: %v", err)
	}

	q, err = queue.GetConnectionToQueue(config.App.Queue.Address)
	if err != nil {
		return fmt.Errorf("error encountered when starting queue service: %v", err)
	}
	defer func() {
		q.Close()
	}()

	if err = server.StartGRPCServer(q, config.App.GRPCServer); err != nil {
		return err
	}
	return nil
}

func StartWorker(workerName string) error {
	var err error
	var q queue.QueueClient
	var db database.DBClient
	var w worker.WorkerService

	err = SetLogOutput(config.App.Logger, "eventhandler.log")
	if err != nil {
		return fmt.Errorf("error encountered when setting log output: %v", err)
	}

	q, err = queue.GetConnectionToQueue(config.App.Queue.Address)
	if err != nil {
		return fmt.Errorf("error encountered when starting queue service: %v", err)
	}

	db, err = database.InitDatabaseConnection(config.App.Database)
	if err != nil {
		return fmt.Errorf("error encountered when connecting to DB: %v", err)
	}

	w, err = worker.CreateWorker(workerName, db, q)
	if err != nil {
		return err
	}

	defer func() {
		db.Close()
		q.Close()
	}()

	log.Printf("Starting worker: %v", workerName)

	w.ProcessWork(workerName)
	return nil
}
