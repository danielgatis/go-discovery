package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/danielgatis/go-discovery"
	"github.com/sirupsen/logrus"
)

func main() {
	d := discovery.NewMdnsDiscovery(5*time.Second, logrus.StandardLogger(), func() (string, string, string, int) {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		return fmt.Sprintf("test:%d", port), "_test._tcp", "local.", port
	})

	output, err := d.Start()
	if err != nil {
		logrus.Fatal(err)
	}

	for peers := range output {
		for i := 0; i < len(peers); i++ {
			peer := peers[i]
			logrus.Info(peer)
		}
	}
}
