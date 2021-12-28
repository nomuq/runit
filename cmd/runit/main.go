package main

import (
	"context"
	"log"
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
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
