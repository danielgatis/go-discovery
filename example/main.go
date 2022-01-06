package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/danielgatis/go-ctrlc"
	"github.com/danielgatis/go-discovery"
	"github.com/sirupsen/logrus"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", 3001, "port number")
}

func main() {
	flag.Parse()

	discovery := discovery.NewMdnsDiscovery(fmt.Sprintf("test:%d", port), "_test._tcp", "local.", port, logrus.StandardLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	ctrlc.Watch(func() {
		cancel()
	})

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
