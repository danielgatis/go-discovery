package discovery

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	m1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// K8sDiscovery is a k8s resolver.
type K8sDiscovery struct {
	clientset kubernetes.Interface
	namespace string
	portName  string
	labels    map[string]string
	logger    logrus.FieldLogger
}

// NewK8sDiscovery returns a new k8s resolver.
func NewK8sDiscovery(clientset kubernetes.Interface, namespace string, portName string, labels map[string]string, logger logrus.FieldLogger) *K8sDiscovery {
	return &K8sDiscovery{
		clientset: clientset,
		namespace: namespace,
		portName:  portName,
		labels:    labels,
		logger:    logger,
	}
}

// Lookup implements discovery.Lookup.
func (d *K8sDiscovery) Lookup() ([]string, error) {
	services, err := d.clientset.CoreV1().Services(d.namespace).List(context.Background(), m1.ListOptions{
		LabelSelector: labels.SelectorFromSet(d.labels).String(),
		Watch:         false,
	})

	if err != nil {
		return []string{}, fmt.Errorf(`d.clientset.CoreV1().Services(d.namespace).List(...): %w`, err)
	}

	peers := make([]string, 0)

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
				peers = append(peers, fmt.Sprintf("%v:%v", podIP, podPort.ContainerPort))
			}
		}
	}

	return peers, nil
}
