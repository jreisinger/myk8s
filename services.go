package main

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetServices(client kubernetes.Clientset, namespace string) (*v1.ServiceList, error) {
	return client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
}

func ToMySvc(in v1.Service, replace, with string) MySvc {
	return MySvc{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: Metadata{
			Name:      strings.ReplaceAll(in.Name, replace, with),
			Namespace: in.Namespace,
			Labels:    in.Labels,
		},
		Spec: Spec{
			Ports:    in.Spec.Ports,
			Selector: in.Spec.Selector,
			Type:     in.Spec.Type,
		},
	}
}

type MySvc struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Spec       `json:"spec,omitempty"`
	Metadata   `json:"metadata"`
}
type Spec struct {
	Ports    []v1.ServicePort  `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"port"`
	Selector map[string]string `json:"selector,omitempty"`
	Type     v1.ServiceType    `json:"type,omitempty"`
}
type Metadata struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}
