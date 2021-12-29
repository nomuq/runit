package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"runit/internal"
	"runit/internal/proto"

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
		Commands: []*cli.Command{
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "debug",
						Aliases: []string{"d"},
						Usage:   "enable debug mode",
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("debug") {
						logrus.SetLevel(logrus.DebugLevel)
					}

					config := proto.DeployRequest{
						Name:       "",
						Repository: "",
						Branch:     "",
						Dockerfile: "",
						Ports: []string{
							"8080",
						},
					}
					json, err := json.MarshalIndent(config, "", "  ")
					if err != nil {
						return err
					}
					// config := &internal.Config{}

					// err := config.Prompt()
					// if err != nil {
					// 	return err
					// }

					// json, err := config.ToJSON()
					// if err != nil {
					// 	return err
					// }

					f, err := os.Create("runit.json")
					if err != nil {
						return err
					}

					_, err = f.Write(json)
					if err != nil {
						return err
					}

					logrus.Infoln("runit.json created")

					return nil
				},
			},

			{
				Name: "deploy",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "debug",
						Aliases: []string{"d"},
						Usage:   "enable debug mode",
						Value:   false,
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "config file path",
						Value:   "runit.json",
					},
				},
				Action: func(c *cli.Context) error {

					// if debug mode is enabled, set log level to debug
					if c.Bool("debug") {
						logrus.SetLevel(logrus.DebugLevel)
					}

					// load config
					// config, err := internal.LoadConfig(c.String("config"))
					// if err != nil {
					// 	return err
					// }

					logrus.Debugln("Connecting to remote host via SSH")
					sshclient := internal.NewSSHClient("localhost:22", "satishbabariya")
					client, err := sshclient.ConnectWithPassword("password")
					if err != nil {
						return err
					}

					dialer := func(ctx context.Context, addr string) (net.Conn, error) {
						// return net.Dial(protocol, addr)
						return client.Dial(internal.Protocol, addr)
					}

					logrus.Debugln("Creating gRPC client")
					conn, err := grpc.Dial(internal.SocketAddress, grpc.WithInsecure(), grpc.WithContextDialer(dialer))
					if err != nil {
						log.Fatal(err)
					}
					defer conn.Close()

					logrus.Debugln("Connected to remote host via gRPC")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
