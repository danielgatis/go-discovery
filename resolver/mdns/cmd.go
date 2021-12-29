package mdns

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/danielgatis/go-discovery/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.AddCmd(NewCmd())
}

// NewCmd returns a new command for mDNS resolver.
func NewCmd() *cobra.Command {
	var (
		instance string
		service  string
		domain   string
		port     int
		interval time.Duration
	)

	var c = &cobra.Command{
		Use: "mdns",
		RunE: func(cmd *cobra.Command, args []string) error {
			kill := make(chan os.Signal, 1)
			signal.Notify(kill, os.Interrupt)

			provider := New(
				instance,
				service,
				port,
				WithDomain(domain),
				WithInterval(interval),
			)

			entries, err := provider.Start()
			if err != nil {
				return fmt.Errorf("provider.Start(...): %w", err)
			}

			go func() {
				for entry := range entries {
					fmt.Println(entry)
				}
			}()

			<-kill
			provider.Stop()
			return nil
		},
	}

	c.Flags().StringVar(&instance, "instance", "", "The instance. Example: my-foo")
	c.MarkFlagRequired("instance")

	c.Flags().StringVar(&service, "service", "", "The service. Example: _foo._tcp")
	c.MarkFlagRequired("service")

	c.Flags().IntVar(&port, "port", 0, "The port. Example: 8000")
	c.MarkFlagRequired("port")

	c.Flags().StringVar(&domain, "domain", "local.", "The domain. Example: local.")
	c.Flags().DurationVar(&interval, "interval", 2*time.Second, "The lookup interval. Example: 2s")

	return c
}
