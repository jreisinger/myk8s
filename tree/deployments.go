package tree

import (
	"fmt"

	"github.com/jreisinger/myk8s/internal/get"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Deployment struct {
	Name        string
	ReplicaSets []ReplicaSet
}
type ReplicaSet struct {
	Name string
	Pods []Pod
}

func PrintDeployments(deployemnts []Deployment) {
	for _, deployment := range deployemnts {
		fmt.Printf("deploy/%s\n", deployment.Name)
		for _, rs := range deployment.ReplicaSets {
			fmt.Printf("%srs/%s\n", "└─", rs.Name)
			for _, pod := range rs.Pods {
				fmt.Printf("%spod/%s\n", "  └─", pod.Name)
				for _, c := range pod.Containers {
					fmt.Printf("%scontainer/%s\n", "    └─", c.Name)
				}
			}
		}
	}
}

func GetDeployments(client kubernetes.Clientset, namespace string) ([]Deployment, error) {
	deploymentList, err := get.Deployments(client, namespace)
	if err != nil {
		return nil, err
	}
	var deployments []Deployment
	for _, deployment := range deploymentList.Items {
		replicaSets, err := getReplicaSetsForDeployment(client, namespace, deployment.Name)
		if err != nil {
			return deployments, err
		}
		var replicasets []ReplicaSet
		for _, rs := range replicaSets {
			var pods []Pod
			v1pods, err := getPodsForReplicaSet(client, namespace, rs.Name)
			if err != nil {
				return deployments, err
			}
			for _, v1pod := range v1pods {
				containers := getPodsContainers(v1pod)
				pods = append(pods, Pod{Name: v1pod.Name, Containers: containers})
			}
			replicasets = append(replicasets, ReplicaSet{Name: rs.Name, Pods: pods})
		}
		deployments = append(deployments, Deployment{Name: deployment.Name, ReplicaSets: replicasets})
	}
	return deployments, nil
}

func getReplicaSetsForDeployment(client kubernetes.Clientset, namespace string, deploymentName string) ([]ReplicaSet, error) {
	replicaSets, err := get.ReplicaSets(client, namespace)
	if err != nil {
		return nil, err
	}
	var deploymentReplicaSets []ReplicaSet
	for _, rs := range replicaSets.Items {
		for _, ownerRef := range rs.OwnerReferences {
			if ownerRef.Kind == "Deployment" && ownerRef.Name == deploymentName {
				deploymentReplicaSets = append(deploymentReplicaSets, ReplicaSet{Name: rs.Name})
			}
		}
	}
	return deploymentReplicaSets, nil
}

func getPodsForReplicaSet(client kubernetes.Clientset, namespace string, replicaSetName string) ([]v1.Pod, error) {
	podList, err := get.Pods(client, namespace, "")
	if err != nil {
		return nil, err
	}
	var pods []v1.Pod
	for _, pod := range podList.Items {
		for _, ownerRef := range pod.OwnerReferences {
			if ownerRef.Kind == "ReplicaSet" && ownerRef.Name == replicaSetName {
				pods = append(pods, pod)
			}
		}
	}
	return pods, nil
}
