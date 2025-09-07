package pulse

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset kubernetes.Interface
}

func NewClient() (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientset,
	}, nil
}

func (c *Client) GetPodStatuses() ([]PodStatus, error) {
	return c.GetPodStatusesInNamespace("")
}

func (c *Client) GetPodStatusesInNamespace(namespace string) ([]PodStatus, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podStatuses []PodStatus
	for _, pod := range pods.Items {
		var restarts int32
		var lastRestart time.Time

		for _, status := range pod.Status.ContainerStatuses {
			restarts += status.RestartCount
			if status.LastTerminationState.Terminated != nil {
				if status.LastTerminationState.Terminated.FinishedAt.After(lastRestart) {
					lastRestart = status.LastTerminationState.Terminated.FinishedAt.Time
				}
			}
		}

		podStatuses = append(podStatuses, PodStatus{
			Name:        pod.Name,
			Namespace:   pod.Namespace,
			Status:      string(pod.Status.Phase),
			Restarts:    restarts,
			LastRestart: lastRestart,
		})
	}

	return podStatuses, nil
}
