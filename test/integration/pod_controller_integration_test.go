//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/quietinvestor/simplecontroller/internal/setup"
	"k8s.io/klog/v2/textlogger"
)

func TestPodLabelReconcilerIntegration(t *testing.T) {
	testEnv := &envtest.Environment{}
	cfg, err := testEnv.Start()
	if err != nil {
		t.Fatalf("failed to start test environment: %v", err)
	}
	t.Cleanup(func() {
		_ = testEnv.Stop()
	})

	loggerConfig := textlogger.NewConfig(textlogger.Verbosity(2))

	mgr, err := setup.Setup(cfg, "default", *loggerConfig)
	if err != nil {
		t.Fatalf("failed to set up manager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	go func() {
		if err := mgr.Start(ctx); err != nil {
			t.Logf("manager exited unexpectedly: %v", err)
			cancel()
		}
	}()

	k8sClient := mgr.GetClient()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"simplecontroller.io/watch": "true",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "test-container",
				Image: "test.internal/integration:latest",
			}},
		},
	}
	if err := k8sClient.Create(ctx, pod); err != nil {
		t.Fatalf("failed to create pod: %v", err)
	}

	updated := &corev1.Pod{}

	// Wait for reconciliation to apply the label
	err = wait.PollUntilContextTimeout(ctx, 100*time.Millisecond, 5*time.Second, true, func(ctx context.Context) (bool, error) {
		if err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-pod", Namespace: "default"}, updated); err != nil {
			return false, nil
		}
		return updated.Labels["simplecontroller.io/color"] == "blue", nil
	})
	if err != nil {
		if err == context.DeadlineExceeded {
			t.Logf("last observed labels: %v", updated.Labels)
			t.Fatal("timed out waiting for label to be applied")
		}
		t.Fatalf("error while polling: %v", err)
	}
}
