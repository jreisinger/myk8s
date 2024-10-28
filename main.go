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

	"github.com/jreisinger/myk8s/dup"
	"github.com/jreisinger/myk8s/internal/clientset"
	"github.com/jreisinger/myk8s/internal/get"
	"github.com/jreisinger/myk8s/logs"
	"github.com/jreisinger/myk8s/tree"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("myk8s: ")

	var (
		kubeconfig string
		namespace  string
	)

	app := &cli.App{
		Usage: "talks to Kubernetes cluster, my way :-)",
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
				Name:      "dup",
				Usage:     "prints existing resources in YAML consumable by kubectl apply",
				ArgsUsage: "kind [name]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "replace",
						Usage: "replace all non-overlapping instances of `string` in name",
					},
					&cli.StringFlag{
						Name:  "with",
						Usage: "with `string`",
					},
				},
				Action: func(cCtx *cli.Context) error {
					supportedKinds := "svc"
					if !cCtx.Args().Present() {
						return fmt.Errorf("please supply one of: %s", supportedKinds)
					}
					client, err := clientset.GetOutOfCluster(kubeconfig)
					if err != nil {
						return err
					}
					switch cCtx.Args().First() {
					case "svc":
						svcs, err := dup.Services(*client, namespace, cCtx.Args().Get(1), cCtx.String("replace"), cCtx.String("with"))
						if err != nil {
							return err
						}
						for _, svc := range svcs {
							fmt.Printf("---\n")
							yamlData, err := yaml.Marshal(&svc)
							if err != nil {
								return err
							}
							fmt.Print(string(yamlData))
						}
					default:
						return fmt.Errorf("unsupported resource kind: %s", cCtx.Args().First())
					}
					return nil
				},
			},
			{
				Name:  "logs",
				Usage: "prints containers logs",
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
					&phase,
				},
				ArgsUsage: "[pod...]",
				Action: func(cCtx *cli.Context) error {
					client, err := clientset.GetOutOfCluster(kubeconfig)
					if err != nil {
						return err
					}
					rx, err := regexp.Compile("(?i)" + cCtx.String("pattern"))
					if err != nil {
						return err
					}
					args := cCtx.Args()
					return logs.Print(client, namespace, rx, cCtx.Int("tail"), cCtx.String("phase"), args.Slice()...)
				},
			},
			{
				Name:  "names",
				Usage: "print object names",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "of",
						Value: "pods",
						Usage: "`kind`",
					},
				},
				Action: func(ctx *cli.Context) error {
					client, err := clientset.GetOutOfCluster(kubeconfig)
					if err != nil {
						return err
					}
					switch ctx.String("of") {
					case "pods":
						podList, err := get.Pods(*client, namespace, "")
						if err != nil {
							return err
						}
						for _, pod := range podList.Items {
							fmt.Println(pod.Name)
						}
					case "deployments":
						deploymentList, err := get.Deployments(*client, namespace)
						if err != nil {
							return err
						}
						for _, d := range deploymentList.Items {
							fmt.Println(d.Name)
						}
					default:
						supportedKinds := []string{"pods", "deployments"}
						return fmt.Errorf("unsupported kind %s, select from: %s", ctx.String("of"), strings.Join(supportedKinds, ", "))
					}
					return nil
				},
			},
			{
				Name:  "tree",
				Usage: "prints top-down relations of objects",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "only-env-to-svc",
						Usage: "only environment variables whose value contains a service name",
					},
				},
				Action: func(cCtx *cli.Context) error {
					client, err := clientset.GetOutOfCluster(kubeconfig)
					if err != nil {
						return err
					}
					services, err := tree.Services(*client, namespace)
					if err != nil {
						return err
					}
					tree.PrintMyServices(services, cCtx.Bool("only-env-to-svc"))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
