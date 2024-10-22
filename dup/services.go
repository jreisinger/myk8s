package dup

import (
	"strings"

	"github.com/jreisinger/myk8s/internal/get"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// Services gets existing service or sevices (if name is empty), removes fields
// not consumable by kubectl apply -f and replaces instances of string in name
// field.
func Services(client kubernetes.Clientset, namespace, name, replace, with string) ([]MyService, error) {
	svcs, err := get.Services(client, namespace)
	if err != nil {
		return nil, err
	}
	var myServices []MyService
	for _, svc := range svcs.Items {
		if name != "" && svc.Name != name {
			continue
		}
		mySvc := toMyService(svc)
		mySvc.Name = strings.ReplaceAll(svc.Name, replace, with)
		myServices = append(myServices, mySvc)
	}
	return myServices, nil
}

func toMyService(in v1.Service) MyService {
	return MyService{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: Metadata{
			Name:        in.Name,
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

type MyService struct {
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
