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
			Name:        strings.ReplaceAll(in.Name, replace, with),
			Namespace:   in.Namespace,
			Labels:      in.Labels,
			Annotations: in.Annotations,
		},
		Spec: Spec{
			Ports:                 in.Spec.Ports,
			Selector:              in.Spec.Selector,
			Type:                  in.Spec.Type,
			InternalTrafficPolicy: in.Spec.InternalTrafficPolicy,
			IPFamilyPolicy:        in.Spec.IPFamilyPolicy,
			IPFamilies:            in.Spec.IPFamilies,
			SessionAffinityConfig: in.Spec.SessionAffinityConfig,
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
	Ports                 []v1.ServicePort                 `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"port"`
	Selector              map[string]string                `json:"selector,omitempty"`
	Type                  v1.ServiceType                   `json:"type,omitempty"`
	InternalTrafficPolicy *v1.ServiceInternalTrafficPolicy `json:"internalTrafficPolicy,omitempty" protobuf:"bytes,22,opt,name=internalTrafficPolicy"`
	IPFamilyPolicy        *v1.IPFamilyPolicy               `json:"ipFamilyPolicy,omitempty" protobuf:"bytes,17,opt,name=ipFamilyPolicy,casttype=IPFamilyPolicy"`
	IPFamilies            []v1.IPFamily                    `json:"ipFamilies,omitempty" protobuf:"bytes,19,opt,name=ipFamilies,casttype=IPFamily"`
	SessionAffinityConfig *v1.SessionAffinityConfig        `json:"sessionAffinityConfig,omitempty" protobuf:"bytes,14,opt,name=sessionAffinityConfig"`
}
type Metadata struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
}
