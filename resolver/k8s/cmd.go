package k8s

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/danielgatis/go-discovery/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func init() {
	cmd.AddCmd(NewCmd())
}

// NewCmd returns a new command for k8s resolver.
func NewCmd() *cobra.Command {
	var (
		namespace string
		portName  string
		labels    map[string]string
		interval  time.Duration
	)

	var c = &cobra.Command{
		Use: "k8s",
		RunE: func(cmd *cobra.Command, args []string) error {
			kill := make(chan os.Signal, 1)
			signal.Notify(kill, os.Interrupt)

			config, err := rest.InClusterConfig()
			if err != nil {
				return fmt.Errorf("rest.InClusterConfig(...): %w", err)
			}

			clientSet, err := kubernetes.NewForConfig(config)
			if err != nil {
				return fmt.Errorf("kubernetes.NewForConfig(...): %w", err)
			}

			provider := New(
				clientSet,
				portName,
				WithNamespace(namespace),
				WithLabels(labels),
				WithInterval(interval),
			)

			entries, err := provider.Start()
			if err != nil {
				logrus.Fatal(err)
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

	c.Flags().StringVar(&portName, "portname", "", "The portName. Example: http")
	c.MarkFlagRequired("portname")

	c.Flags().StringVar(&namespace, "namespace", "default", "The namespace. Example: default")
	c.Flags().StringToStringVar(&labels, "labels", make(map[string]string), "the labels. Example: foo=bar")
	c.Flags().DurationVar(&interval, "interval", 2*time.Second, "The lookup interval. Example: 2s")

	return c
}
