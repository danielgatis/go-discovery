package k8s

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	m1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Discovery is a k8s resolver.
type Discovery struct {
	clientset kubernetes.Interface
	namespace string
	portName  string
	labels    map[string]string
	logger    logrus.FieldLogger
	interval  time.Duration
	output    chan string
	stop      chan struct{}
}

// New returns a new k8s resolver.
func New(clientset kubernetes.Interface, portName string, opts ...Option) *Discovery {
	const (
		defaultNamespace = "default"
		defaultInterval  = 5 * time.Second
	)

	var (
		defaultLabels = make(map[string]string)
		defaultLogger = logrus.StandardLogger()
	)

	d := &Discovery{
		clientset: clientset,
		namespace: defaultNamespace,
		portName:  portName,
		labels:    defaultLabels,
		logger:    defaultLogger,
		interval:  defaultInterval,
		output:    make(chan string),
		stop:      make(chan struct{}),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

// Start implements resolver.Resolver.
func (d *Discovery) Start() (chan string, error) {
	ticker := time.NewTicker(d.interval)

	go func() {
		for {
			select {
			case <-d.stop:
				ticker.Stop()
				return
			case <-ticker.C:
				services, err := d.clientset.CoreV1().Services(d.namespace).List(context.Background(), m1.ListOptions{
					LabelSelector: labels.SelectorFromSet(d.labels).String(),
					Watch:         false,
				})

				if err != nil {
					d.logger.Errorf("Error during k8s service lookup: %v\n", err)
					continue
				}

				for _, service := range services.Items {
					pods, err := d.clientset.CoreV1().Pods(service.Namespace).List(context.Background(), m1.ListOptions{
						LabelSelector: labels.SelectorFromSet(labels.Set(service.Spec.Selector)).String(),
					})

					if err != nil {
						d.logger.Errorf("Error during k8s pod lookup: %v\n", err)
						continue
					}

					for _, pod := range pods.Items {
						if strings.ToLower(string(pod.Status.Phase)) != "running" {
							continue
						}

						podIP := pod.Status.PodIP
						var podPort v1.ContainerPort

						for _, container := range pod.Spec.Containers {
							for _, port := range container.Ports {
								if port.Name == d.portName {
									podPort = port
									break
								}
							}
						}

						if podIP != "" && podPort.ContainerPort != 0 {
							d.output <- fmt.Sprintf("%v:%v", podIP, podPort.ContainerPort)
						}
					}
				}
			}
		}
	}()

	return d.output, nil
}

// Stop implements resolver.Resolver.
func (d *Discovery) Stop() {
	d.stop <- struct{}{}
	close(d.output)
}
