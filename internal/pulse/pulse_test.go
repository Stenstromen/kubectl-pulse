package pulse

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetClusterPulseHealthy(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	service, err := NewServiceWithClientset(clientset)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	result, err := service.GetClusterPulse(15, 30)
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	t.Log(result)
}

func TestGetClusterPulseWarning(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	fakePods := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-1",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							RestartCount: 1,
							LastTerminationState: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									FinishedAt: metav1.NewTime(time.Now().Add(-5 * time.Minute)),
								},
							},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-2",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							RestartCount: 3,
							LastTerminationState: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									FinishedAt: metav1.NewTime(time.Now().Add(-2 * time.Minute)),
								},
							},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-3",
					Namespace: "kube-system",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodFailed,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							RestartCount: 0,
						},
					},
				},
			},
		},
	}

	for _, pod := range fakePods.Items {
		_, err := clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Failed to create fake pod: %v", err)
		}
	}

	service, err := NewServiceWithClientset(clientset)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	result, err := service.GetClusterPulse(15, 3)
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	t.Log(result)
}
