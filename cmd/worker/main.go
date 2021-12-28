package main

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const (
	protocol = "unix"
	sockAddr = "/tmp/runit.sock"
)

func main() {
	app := &cli.App{
		Name:  "runit",
		Usage: "run docker containers on remote hosts with ease",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			logrus.Println("boom! I say!")

			// crate unix socket server

			if _, err := os.Stat(sockAddr); err == nil {
				if err := os.RemoveAll(sockAddr); err != nil {
					return err
				}
			}

			listener, err := net.Listen(protocol, sockAddr)
			if err != nil {
				return err
			}

			server := grpc.NewServer()

			return server.Serve(listener)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		logrus.Println(string(buf[:n]))
	}
}
