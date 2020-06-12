package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	clientApp := cli.NewApp()
	clientApp.Name = "Squadcast service"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "start",
			Description: "Start Http Server",
			Action: func(c *cli.Context) error {
				fmt.Println("Let's begin")
				return nil
			},
		},
	}
	if err := clientApp.Run(os.Args); err != nil {
		panic(err)
	}
}
