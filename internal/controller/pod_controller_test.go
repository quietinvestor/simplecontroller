package controller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestPodLabelReconcilerReconcile(t *testing.T) {
	tests := []struct {
		name        string
		pod         *corev1.Pod
		expectPatch bool
	}{
		{
			name:        "Pod not found",
			pod:         nil,
			expectPatch: false,
		},
		{
			name:        "Missing watch label",
			pod:         &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "default"}},
			expectPatch: false,
		},
		{
			name: "Watch label not true",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod2", Namespace: "default",
					Labels: map[string]string{WatchKey: "no"},
				},
			},
			expectPatch: false,
		},
		{
			name: "Already has color label",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod3", Namespace: "default",
					Labels: map[string]string{WatchKey: "true", ColorKey: "blue"},
				},
			},
			expectPatch: false,
		},
		{
			name: "Should patch label",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod4", Namespace: "default",
					Labels: map[string]string{WatchKey: "true"},
				},
			},
			expectPatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			if err := corev1.AddToScheme(scheme); err != nil {
				t.Fatalf("failed to add corev1 to scheme: %v", err)
			}

			var fakeClient client.Client
			if tt.pod != nil {
				fakeClient = fake.NewClientBuilder().WithScheme(scheme).WithObjects(tt.pod).Build()
			} else {
				fakeClient = fake.NewClientBuilder().WithScheme(scheme).Build()
			}

			r := &PodLabelReconciler{Client: fakeClient}
			ctx := context.Background()

			// Reconcile always receives a NamespacedName, even if the Pod no longer exists.
			// Use a fallback name to simulate the case where the Pod was deleted after being queued.
			req := types.NamespacedName{Name: "deletedpod", Namespace: "default"}
			if tt.pod != nil {
				req = types.NamespacedName{Name: tt.pod.Name, Namespace: tt.pod.Namespace}
			}

			before := ""
			if tt.pod != nil {
				before = tt.pod.Labels[ColorKey]
			}

			_, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: req})
			if err != nil {
				t.Fatalf("Reconcile returned error: %v", err)
			}

			if tt.pod != nil {
				var pod corev1.Pod

				err := fakeClient.Get(ctx, req, &pod)
				if err != nil {
					t.Fatalf("failed to get pod: %v", err)
				}

				after := pod.Labels[ColorKey]

				if tt.expectPatch && after != ColorValue {
					t.Fatalf("expected color label to be %q, got %q", ColorValue, after)
				}
				if !tt.expectPatch && after != before {
					t.Fatalf("expected color label to remain %q, but got %q", before, after)
				}
			}
		})
	}
}
