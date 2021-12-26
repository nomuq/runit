package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "runit",
		Usage: "run docker containers on remote hosts with ease",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			logrus.Println("boom! I say!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
