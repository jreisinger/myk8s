package logs

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jreisinger/myk8s/get"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Container struct {
	Name string
	Logs []string
}

func logs(client kubernetes.Clientset, namespace string, pod corev1.Pod, regex *regexp.Regexp) ([]Container, error) {
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

func Print(client *kubernetes.Clientset, namespace string, rx *regexp.Regexp, tail int, podPhase string, podNames ...string) error {
	pods, err := get.Pods(*client, namespace, podPhase)
	if err != nil {
		return err
	}

	podHasLogs := make(map[string]bool)
POD:
	for _, pod := range pods.Items {
		if !found(pod.Name, podNames...) {
			continue POD
		}

		containers, _ := logs(*client, namespace, pod, rx)
		for _, c := range containers {
			if len(c.Logs) > 0 {
				podHasLogs[pod.Name] = true
				continue POD
			}
		}

	}

	for _, pod := range pods.Items {
		if !podHasLogs[pod.Name] {
			continue
		}

		containers, err := logs(*client, namespace, pod, rx)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("=== %s ===\n", pod.Name)
		for _, c := range containers {
			fmt.Printf("--- %s ---\n", c.Name)
			for _, log := range lastNElements(c.Logs, tail) {
				fmt.Println(log)
			}
		}
	}

	return nil
}

func found(name string, names ...string) bool {
	if len(names) == 0 {
		return true
	}
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func lastNElements(ss []string, n int) []string {
	if n > len(ss) {
		return ss
	}
	return ss[len(ss)-n:]
}
