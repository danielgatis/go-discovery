package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgatis/go-discovery"
	"github.com/sirupsen/logrus"
)

var (
	port *int
)

func init() {
	port = flag.Int("port", 3001, "port number")
}

func main() {
	flag.Parse()
	discovery := discovery.NewMdnsDiscovery(fmt.Sprintf("test:%d", port), "_test._tcp", "local.", *port, logrus.StandardLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exit
		cancel()
	}()

	go func() {
		discovery.Register(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			peers, err := discovery.Lookup()
			if err != nil {
				logrus.Fatal(err)
			}

			for _, peer := range peers {
				logrus.Info(peer)
			}
		}
	}
}
