```
NAME:
   myk8s - talks to Kubernetes cluster, my way :-)

USAGE:
   myk8s [global options] command [command options]

COMMANDS:
   dup      prints existing resources in YAML consumable by kubectl apply
   graph    prints top-down relations of a resource kind
   logs     prints containers logs
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --kubeconfig value, -k value  path to the kubeconfig file (default: "/Users/jozef.reisinger/.kube/config")
   --namespace value, -n value   kubernetes namespace (default: "default")
   --help, -h                    show help
```