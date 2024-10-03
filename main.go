package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"
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
				Usage: "Prints logs of containers",
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
					return Logs(client, namespace, rx, cCtx.Int("tail"), cCtx.String("phase"), args.Slice()...)
				},
			},
			{
				Name:  "services",
				Usage: "Prints services in YAML with useless fields removed",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "replace",
						Usage: "replace all non-overlapping instances of `string` in service name",
					},
					&cli.StringFlag{
						Name:  "with",
						Usage: "with `string`",
					},
				},
				Action: func(cCtx *cli.Context) error {
					client, err := GetOutOfClusterClient(kubeconfig)
					if err != nil {
						return err
					}
					services, err := GetServices(*client, namespace)
					if err != nil {
						return err
					}
					for _, svc := range services.Items {
						mySvc := ToMySvc(svc, cCtx.String("replace"), cCtx.String("with"))
						fmt.Printf("---\n")
						yamlData, err := yaml.Marshal(&mySvc)
						if err != nil {
							return err
						}
						fmt.Print(string(yamlData))
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
