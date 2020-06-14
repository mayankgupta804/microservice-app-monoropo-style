package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/squadcast_assignment/internal/config"
	"github.com/squadcast_assignment/internal/infrastructure/database"
	"github.com/squadcast_assignment/internal/migrations"
	"github.com/squadcast_assignment/internal/repository"
	"github.com/squadcast_assignment/internal/service"
	"github.com/squadcast_assignment/internal/webserver"

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
				if err := StartIncidentWebService(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:        "start:eventhandler",
			Description: "Start Event Handler Service",
			Action: func(c *cli.Context) error {
				return errors.New("Not implemented")
			},
		},
		{
			Name:        "start:worker",
			Description: "Start worker of given type",
			Action: func(c *cli.Context) error {
				return errors.New("Not implemented")
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

	db, err = database.InitDatabaseConnection(config.App.Database)
	if err != nil {
		return fmt.Errorf("error encountered when connecting to DB: %v", err)
	}

	repo = repository.InitIncidentRepository(db)
	incidentService = service.NewIncidentService(repo)
	router := webserver.SetupRoutes(incidentService)
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
