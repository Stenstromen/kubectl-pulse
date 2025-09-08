package pulse

import (
	"context"
	"strings"
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

	result, err := service.GetClusterPulse(15, 30, "")
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

	result, err := service.GetClusterPulse(15, 3, "")
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	t.Log(result)
}

func TestGetClusterPulseWithNamespace(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	fakePods := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-default-1",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							RestartCount: 2,
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
					Name:      "pod-kube-system-1",
					Namespace: "kube-system",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							RestartCount: 5,
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
					Name:      "pod-kube-system-2",
					Namespace: "kube-system",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							RestartCount: 1,
							LastTerminationState: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									FinishedAt: metav1.NewTime(time.Now().Add(-10 * time.Minute)),
								},
							},
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

	// Test with kube-system namespace
	result, err := service.GetClusterPulse(15, 3, "kube-system")
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	t.Logf("kube-system namespace result: %s", result)

	// Test with default namespace
	result, err = service.GetClusterPulse(15, 3, "default")
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	t.Logf("default namespace result: %s", result)

	// Test with all namespaces (empty string)
	result, err = service.GetClusterPulse(15, 3, "")
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	t.Logf("all namespaces result: %s", result)
}

func TestPodStatusDistribution(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	// Create pods with different statuses
	fakePods := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "running-pod-1",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{RestartCount: 0},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "running-pod-2",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{RestartCount: 1},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pending-pod",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodPending,
					ContainerStatuses: []corev1.ContainerStatus{
						{RestartCount: 0},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "failed-pod",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodFailed,
					ContainerStatuses: []corev1.ContainerStatus{
						{RestartCount: 0},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "succeeded-pod",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodSucceeded,
					ContainerStatuses: []corev1.ContainerStatus{
						{RestartCount: 0},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "unknown-pod",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodUnknown,
					ContainerStatuses: []corev1.ContainerStatus{
						{RestartCount: 0},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "other-namespace-pod",
					Namespace: "kube-system",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{RestartCount: 0},
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

	// Test all namespaces
	result, err := service.GetClusterPulse(15, 10, "")
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	// Verify the output contains pod status distribution
	if !strings.Contains(result, "📊 Pod Status Distribution:") {
		t.Error("Expected output to contain pod status distribution section")
	}

	// Test specific namespace filtering
	result, err = service.GetClusterPulse(15, 10, "default")
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	// Verify the output contains pod status distribution for default namespace only
	if !strings.Contains(result, "📊 Pod Status Distribution:") {
		t.Error("Expected output to contain pod status distribution section")
	}

	t.Logf("Pod status distribution test result: %s", result)
}

func TestPodStatusDistributionPercentages(t *testing.T) {
	// Test percentage calculation directly
	distribution := PodStatusDistribution{
		Running:   3,
		Pending:   1,
		Failed:    1,
		Succeeded: 0,
		Unknown:   0,
		Total:     5,
	}

	// Test percentage calculations
	expectedRunning := 60.0
	expectedPending := 20.0
	expectedFailed := 20.0
	expectedSucceeded := 0.0
	expectedUnknown := 0.0

	if got := distribution.GetPercentage("Running"); got != expectedRunning {
		t.Errorf("GetPercentage('Running') = %.1f, want %.1f", got, expectedRunning)
	}

	if got := distribution.GetPercentage("Pending"); got != expectedPending {
		t.Errorf("GetPercentage('Pending') = %.1f, want %.1f", got, expectedPending)
	}

	if got := distribution.GetPercentage("Failed"); got != expectedFailed {
		t.Errorf("GetPercentage('Failed') = %.1f, want %.1f", got, expectedFailed)
	}

	if got := distribution.GetPercentage("Succeeded"); got != expectedSucceeded {
		t.Errorf("GetPercentage('Succeeded') = %.1f, want %.1f", got, expectedSucceeded)
	}

	if got := distribution.GetPercentage("Unknown"); got != expectedUnknown {
		t.Errorf("GetPercentage('Unknown') = %.1f, want %.1f", got, expectedUnknown)
	}

	// Test edge case: zero total
	emptyDistribution := PodStatusDistribution{Total: 0}
	if got := emptyDistribution.GetPercentage("Running"); got != 0.0 {
		t.Errorf("GetPercentage('Running') with zero total = %.1f, want 0.0", got)
	}
}

func TestPodStatusDistributionEdgeCases(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	service, err := NewServiceWithClientset(clientset)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Test with no pods
	result, err := service.GetClusterPulse(15, 10, "")
	if err != nil {
		t.Fatalf("Failed to get cluster pulse: %v", err)
	}

	// Should show "No pods found" message
	if !strings.Contains(result, "📊 Pod Status: No pods found") {
		t.Error("Expected output to show 'No pods found' message when no pods exist")
	}

	t.Logf("Empty cluster test result: %s", result)
}
