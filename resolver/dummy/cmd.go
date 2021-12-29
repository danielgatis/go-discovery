package dummy

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

// NewCmd returns a new command for dummy resolver.
func NewCmd() *cobra.Command {
	var (
		peers    []string
		interval time.Duration
	)

	var c = &cobra.Command{
		Use: "dummy",
		RunE: func(cmd *cobra.Command, args []string) error {
			kill := make(chan os.Signal, 1)
			signal.Notify(kill, os.Interrupt)

			provider := New(
				peers,
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

	c.Flags().StringArrayVar(&peers, "peer", []string{}, "The peer addr as ip:port. Example: 127.0.0.1:80")
	c.MarkFlagRequired("peer")

	c.Flags().DurationVar(&interval, "interval", 2*time.Second, "The lookup interval. Example: 2s")

	return c
}
