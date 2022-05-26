package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// NOTE: RBAC not needed here.
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=Ignore,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io,admissionReviewVersions=v1,sideEffects=NoneOnDryRun

// PodMutator annotates Pods
type PodMutator struct {
	Client  client.Client
	decoder *admission.Decoder
	Log     logr.Logger
}

// PodMutator adds an annotation to every incoming pods.
func (m *PodMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	fmt.Printf("Handling pod %s/%s", req.Namespace, req.Name)
	err := m.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.ObjectMeta.Namespace == "open-feature-operator-system" {
		m.Log.Info("Skipping pod %s/%s", req.Namespace, req.Name)
		return admission.Response{}
	}

	configName := fmt.Sprintf("%s-%s-config", pod.Name, pod.Namespace)
	// Create the agent config
	fmt.Printf("Creating configmap %s/%s", pod.Namespace, configName)
	if err := m.Client.Create(ctx, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: pod.Namespace,
		},
		//TODO
		Data: map[string]string{
			"config.yaml": "{}",
		},
	}); err != nil {
		fmt.Printf("failed to create config map %s", configName)
		return admission.Errored(http.StatusInternalServerError, err)
	}

	fmt.Printf("Creating sidecar for pod %s/%s", pod.Namespace, pod.Name)
	// Inject the agent
	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: "agent-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configName,
				},
			},
		},
	})
	pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
		Name:  "agent",
		Image: "tibbar/of-agent:v0.0.1",
		Args: []string{
			"start", "-f", "/etc/of-agent/config.yaml",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "agent-config",
				MountPath: "/etc/of-agent",
			},
		},
	})

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// PodMutator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (m *PodMutator) InjectDecoder(d *admission.Decoder) error {
	m.decoder = d
	return nil
}
