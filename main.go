package main

import (
	"fmt"
	"net"
	"net/http"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	// Hosts
	hosts := map[string]*Host{}

	// Reverse proxy
	hosts["essentiel.dev"] = ReverseProxy("essentiel.dev", "https://satishbabariya.github.io/")
	hosts["localhost"] = ReverseProxy("localhost", "http://localhost:8080")

	e := echo.New()
	e.HideBanner = true

	for host, _ := range hosts {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(host)
	}

	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()

		host := hosts[req.Host]
		e.Logger.Info(req.Host)

		if host == nil {
			return echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return err
	})
	go func() {
		e.Start(":80")
	}()

	go func() {
		e.StartAutoTLS(":443")
	}()
	// e.Logger.Fatal(e.Start(":80"))
	// e.Logger.Fatal(e.StartAutoTLS(":443"))

	e.Logger.Fatal(Serve())
}

func ReverseProxy(host string, proxy string) *Host {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Any("/*", func(c echo.Context) (err error) {
		c.Redirect(http.StatusFound, proxy)
		return
	})

	return &Host{
		Echo: e,
	}
}

func Serve() error {
	socketpath := "/tmp/runit"
	// carry on with your socket creation:
	addr, err := net.ResolveUnixAddr("unixgram", socketpath)
	if err != nil {
		return err
	}

	// always remove the named socket from the fs if its there
	err = syscall.Unlink(socketpath)
	if err != nil {
		// not really important if it fails
		return err
	}

	// carry on with socket bind()
	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		return err
	}

	// now we can start serving
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}

		fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))
	}
}

// func listen(end chan<- bool) {
// 	addr, err := net.ResolveUnixAddr("unix", "/tmp/runit")
// 	if err != nil {
// 		fmt.Printf("Failed to resolve: %v\n", err)
// 		os.Exit(1)
// 	}

// 	list, err := net.ListenUnix("unix", addr)
// 	if err != nil {
// 		fmt.Printf("failed to listen: %v\n", err)
// 		os.Exit(1)
// 	}
// 	conn, _ := list.AcceptUnix()

// 	buf := make([]byte, 2048)
// 	n, uaddr, err := conn.ReadFromUnix(buf)
// 	if err != nil {
// 		fmt.Printf("LISTEN: Error: %v\n", err)
// 	} else {
// 		fmt.Printf("LISTEN: received %v bytes from %+v\n", n, uaddr)
// 		fmt.Printf("LISTEN: %v\n", string(buf))
// 	}

// 	conn.Close()
// 	list.Close()
// 	end <- true
// }
