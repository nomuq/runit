package main

import (
	"net"
	"os"
	"runit/internal"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func main() {
	app := &cli.App{
		Name:  "runit",
		Usage: "run docker containers on remote hosts with ease",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable debug mode",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			// if debug mode is enabled, set log level to debug
			if c.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}

			logrus.Debugln("Remove socket file")
			if _, err := os.Stat(internal.SocketAddress); err == nil {
				if err := os.RemoveAll(internal.SocketAddress); err != nil {
					return err
				}
			}

			logrus.Debugln("Creating socket listener")
			listener, err := net.Listen(internal.Protocol, internal.SocketAddress)
			if err != nil {
				return err
			}

			logrus.Debugln("Starting gRPC server")
			server := grpc.NewServer()

			logrus.Debugln("gRPC Server is listening on", internal.SocketAddress)
			return server.Serve(listener)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
