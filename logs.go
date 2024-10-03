package main

import (
	"context"
	"regexp"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Container struct {
	Name string
	Logs []string
}

func GetLogs(client kubernetes.Clientset, namespace string, pod corev1.Pod, regex *regexp.Regexp) ([]Container, error) {
	var containers []Container

	for _, container := range pod.Spec.Containers {
		logs, err := client.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name}).Do(context.TODO()).Raw()
		if err != nil {
			return nil, err
		}

		var logsFiltered []string
		for _, line := range strings.Split(string(logs), "\n") {
			if line == "" {
				continue
			}
			if regex.MatchString(line) {
				logsFiltered = append(logsFiltered, line)
			}
		}

		containers = append(containers, Container{
			Name: container.Name,
			Logs: logsFiltered,
		})
	}

	return containers, nil
}
