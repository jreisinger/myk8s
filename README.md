```
NAME:
   myk8s - my way of talking to Kubernetes cluster :-)

USAGE:
   myk8s [global options] command [command options]

COMMANDS:
   logs      Prints logs of containers
   services  Prints services in YAML consumable by kubectl
   graph     Visualizes relations of a resource
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --kubeconfig value, -k value  path to the kubeconfig file (default: "/Users/jozef.reisinger/.kube/config")
   --namespace value, -n value   kubernetes namespace (default: "default")
   --help, -h                    show help
```