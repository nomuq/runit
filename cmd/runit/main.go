package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
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

			sshclient, err := SSHConnectWithPassword("localhost:22", "satishbabariya", "password")
			if err != nil {
				return err
			}

			dialer := func(ctx context.Context, addr string) (net.Conn, error) {
				// return net.Dial(protocol, addr)
				return sshclient.Dial(protocol, addr)
			}

			conn, err := grpc.Dial(sockAddr, grpc.WithInsecure(), grpc.WithContextDialer(dialer))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			logrus.Println("boom! I say!")

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func SSHConnectWithKey(addr, user, key string) (*ssh.Client, error) {
	//signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte("password"))
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func SSHConnectWithPassword(addr, user, password string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
