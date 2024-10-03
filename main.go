package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("myk8s: ")

	var (
		kubeconfig string
		namespace  string
	)

	app := &cli.App{
		Usage: "Talk to Kubernetes cluster, my way :-)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "kubeconfig",
				Aliases:     []string{"k"},
				Value:       filepath.Join(homedir.HomeDir(), ".kube", "config"),
				Usage:       "path to the kubeconfig file",
				Destination: &kubeconfig,
			},
			&cli.StringFlag{
				Name:        "namespace",
				Aliases:     []string{"n"},
				Value:       "default",
				Usage:       "kubernetes namespace",
				Destination: &namespace,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "logs",
				Usage: "Prints containers' logs",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "pattern",
						Usage: "logs matching `regex` pattern (case insensitive)",
					},
					&cli.IntFlag{
						Name:  "tail",
						Value: 10,
						Usage: "last `n` log lines",
					},
					&cli.StringFlag{
						Name:  "phase",
						Usage: "logs of pods in lifecycle `phase`",
						Action: func(ctx *cli.Context, s string) error {
							possible := []string{"pending", "running", "succeeded", "failed", "unknown"}
							for _, p := range possible {
								if ctx.String("phase") == p {
									return nil
								}
							}
							return fmt.Errorf("possible phases: %v", strings.Join(possible, ", "))
						},
					},
				},
				ArgsUsage: "[pod...]",
				Action: func(cCtx *cli.Context) error {
					client, err := GetOutOfClusterClient(kubeconfig)
					if err != nil {
						return err
					}
					rx, err := regexp.Compile("(?i)" + cCtx.String("pattern"))
					if err != nil {
						return err
					}
					args := cCtx.Args()
					return logs(client, namespace, rx, cCtx.Int("tail"), cCtx.String("phase"), args.Slice()...)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func logs(client *kubernetes.Clientset, namespace string, rx *regexp.Regexp, tail int, podPhase string, podNames ...string) error {
	pods, err := GetPods(*client, namespace, podPhase)
	if err != nil {
		return err
	}

	podHasLogs := make(map[string]bool)
POD:
	for _, pod := range pods.Items {
		if !found(pod.Name, podNames...) {
			continue POD
		}

		containers, _ := GetLogs(*client, namespace, pod, rx)
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

		containers, err := GetLogs(*client, namespace, pod, rx)
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
