package graph

import (
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/jreisinger/myk8s/get"
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
}

func PrintMyServices(myServices []MyService) {
	for _, svc := range myServices {
		fmt.Printf("%s: %s\n", svc.Name, formatServicePorts(svc.Ports))
		for _, pod := range svc.Pods {
			fmt.Printf("%s%s\n", "└─", pod.Name)
			for _, c := range pod.Containers {
				fmt.Printf("%s%s: %s\n", "  └─", c.Name, formatContainerPorts(c.Ports))
			}
		}
	}
}

func formatServicePorts(ports []v1.ServicePort) string {
	var ss []string
	for _, p := range ports {
		ss = append(ss, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
	}
	return strings.Join(ss, ",")
}

func formatContainerPorts(ports []v1.ContainerPort) string {
	var ss []string
	for _, p := range ports {
		ss = append(ss, fmt.Sprintf("%d/%s", p.ContainerPort, p.Protocol))
	}
	return strings.Join(ss, ",")
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
		myContainers = append(myContainers, MyContainer{Name: c.Name, Ports: c.Ports})
	}
	return myContainers
}

func isServiceMatchingPod(serviceSelector, podLabels map[string]string) bool {
	for selectorK, selectorV := range serviceSelector {
		if podLabelV, ok := podLabels[selectorK]; !ok || podLabelV != selectorV {
			return false
		}
	}
	return true
}
