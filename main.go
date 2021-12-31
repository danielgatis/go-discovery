package main

import (
	"github.com/danielgatis/go-discovery/cmd"

	_ "github.com/danielgatis/go-discovery/discovery/dummy"
	_ "github.com/danielgatis/go-discovery/discovery/k8s"
	_ "github.com/danielgatis/go-discovery/discovery/mdns"
)

func main() {
	cmd.Execute()
}
