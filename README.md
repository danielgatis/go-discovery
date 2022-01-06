# Go - Discovery

[![Go Report Card](https://goreportcard.com/badge/github.com/danielgatis/go-discovery?style=flat-square)](https://goreportcard.com/report/github.com/danielgatis/go-discovery)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/danielgatis/go-discovery/master/LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/danielgatis/go-discovery)

A collection of service discovery implementations.

## Install

```bash
go get -u github.com/danielgatis/go-discovery
```

And then import the package in your code:

```go
import "github.com/danielgatis/go-discovery"
```

### Example

An example described below is one of the use cases.

```go
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
```

```
❯ go run main.go -p 3001
❯ go run main.go -p 3002
```

### License

Copyright (c) 2021-present [Daniel Gatis](https://github.com/danielgatis)

Licensed under [MIT License](./LICENSE)

### Buy me a coffee

Liked some of my work? Buy me a coffee (or more likely a beer)

<a href="https://www.buymeacoffee.com/danielgatis" target="_blank"><img src="https://bmc-cdn.nyc3.digitaloceanspaces.com/BMC-button-images/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;"></a>
