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
	ctxReg, cancelReg := context.WithCancel(context.Background())
	ctxLkp, cancelLkp := context.WithTimeout(context.Background(), 1*time.Second)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancelReg()
		cancelLkp()
	}()

	go func() {
		discovery.Register(ctxReg)
	}()

	for {
		select {
		case <-ctxReg.Done():
			return
		default:
			ctxLkp, cancelLkp = context.WithTimeout(context.Background(), 1*time.Second)
			peers, err := discovery.Lookup(ctxLkp)
			if err != nil {
				logrus.Fatal(err)
			}

			for _, peer := range peers {
				logrus.Info(peer)
			}
		}
	}
}
