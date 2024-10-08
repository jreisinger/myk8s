package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

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
		Usage: "my way of talking to Kubernetes cluster :-)",
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
					&phase,
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
				Usage: "Prints services in YAML consumable by kubectl",
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
			{
				Name:  "graph",
				Usage: "Visualizes relations of each pod",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "bypod",
						Usage: "group by pod instead by top level object",
					},
					&phase,
				},
				Action: func(cCtx *cli.Context) error {
					client, err := GetOutOfClusterClient(kubeconfig)
					if err != nil {
						return err
					}
					graphs, err := GetGraphs(*client, namespace, cCtx.String("phase"))
					if err != nil {
						return err
					}
					if cCtx.Bool("nogroup") {
						for _, g := range graphs {
							fmt.Println(g)
						}
						return nil
					}
					fmt.Println(graphs)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
