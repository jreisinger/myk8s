package tree

import (
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/jreisinger/myk8s/internal/get"
)

type MyService struct {
	Name  string
	Ports []v1.ServicePort
	Pods  []MyPod
}

type MyPod struct {
	Name       string
	Containers []MyContainer
}

type MyContainer struct {
	Name  string
	Ports []v1.ContainerPort
	Env   []v1.EnvVar
}

func PrintMyServices(myServices []MyService, onlyEnvToSvc bool) {
	for _, svc := range myServices {
		fmt.Printf("svc/%s: %s\n", svc.Name, formatServicePorts(svc.Ports))
		for _, pod := range svc.Pods {
			fmt.Printf("%spod/%s\n", "└─", pod.Name)
			for _, c := range pod.Containers {
				fmt.Printf("%scontainer/%s: %s\n", "  └─", c.Name, formatContainerPorts(c.Ports))
				for _, e := range c.Env {
					if onlyEnvToSvc {
						if s := envVarReferencesService(e, myServices); s != nil {
							fmt.Printf("%s%s: %s", "    └─", e.Name, e.Value)
							fmt.Printf(" -> svc/%s\n", s.Name)
						}
					} else {
						fmt.Printf("%s%s: %s", "    └─", e.Name, e.Value)
						if s := envVarReferencesService(e, myServices); s != nil {
							fmt.Printf(" -> svc/%s", s.Name)
						}
						fmt.Println()
					}
				}
			}
		}
	}
}

func envVarReferencesService(e v1.EnvVar, myServices []MyService) *MyService {
	var longest MyService
	for _, svc := range myServices {
		if strings.Contains(e.Value, svc.Name) {
			if len(svc.Name) > len(longest.Name) {
				longest = svc
			}
		}
	}
	if longest.Name != "" {
		return &longest
	}
	return nil
}

func formatServicePorts(ports []v1.ServicePort) string {
	var ss []string
	for _, p := range ports {
		ss = append(ss, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
	}
	return strings.Join(ss, ", ")
}

func formatContainerPorts(ports []v1.ContainerPort) string {
	var ss []string
	for _, p := range ports {
		ss = append(ss, fmt.Sprintf("%d/%s", p.ContainerPort, p.Protocol))
	}
	return strings.Join(ss, ", ")
}

func Services(client kubernetes.Clientset, namespace string) ([]MyService, error) {
	services, err := get.Services(client, namespace)
	if err != nil {
		return nil, fmt.Errorf("getting services: %v", err)
	}
	pods, err := get.Pods(client, namespace, "")
	if err != nil {
		return nil, fmt.Errorf("getting pods: %v", err)
	}
	var myServices []MyService
	for _, svc := range services.Items {
		if svc.Spec.ClusterIP == "None" && svc.Spec.Selector == nil {
			log.Printf("skipping headless service w/o selector: %s\n", svc.Name)
			continue
		}
		var servicePods []MyPod
		for _, pod := range pods.Items {
			podLabels := pod.GetLabels()
			if isServiceMatchingPod(svc.Spec.Selector, podLabels) {
				servicePods = append(servicePods,
					MyPod{Name: pod.Name, Containers: getPodsContainers(pod)},
				)
			}
		}
		myServices = append(myServices, MyService{
			Name:  svc.Name,
			Ports: svc.Spec.Ports,
			Pods:  servicePods,
		})
	}
	return myServices, nil
}

func getPodsContainers(pod v1.Pod) []MyContainer {
	var myContainers []MyContainer
	for _, c := range pod.Spec.Containers {
		myContainers = append(myContainers, MyContainer{Name: c.Name, Ports: c.Ports, Env: c.Env})
	}
	return myContainers
}

func isServiceMatchingPod(serviceSelector, podLabels map[string]string) bool {
	if serviceSelector == nil {
		return false
	}
	for selectorK, selectorV := range serviceSelector {
		if podLabelV, ok := podLabels[selectorK]; !ok || podLabelV != selectorV {
			return false
		}
	}
	return true
}
