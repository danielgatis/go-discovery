package main

import (
	"github.com/danielgatis/go-discovery/cmd"

	_ "github.com/danielgatis/go-discovery/resolver/dummy"
	_ "github.com/danielgatis/go-discovery/resolver/k8s"
	_ "github.com/danielgatis/go-discovery/resolver/mdns"
)

func main() {
	cmd.Execute()
}
