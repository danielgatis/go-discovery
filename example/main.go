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
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	d := discovery.NewMdnsDiscovery(fmt.Sprintf("test:%d", port), "_test._tcp", "local.", port, 5*time.Second, logrus.StandardLogger())

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
