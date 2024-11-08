```
❯ myk8s 
NAME:
   myk8s - talks to Kubernetes cluster, my way :-)

USAGE:
   myk8s [global options] command [command options]

COMMANDS:
   dup      prints existing resources in YAML consumable by kubectl apply
   logs     prints containers logs
   names    prints object names
   tree     prints top-down relations of objects
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --kubeconfig value, -k value  path to the kubeconfig file (default: "/Users/jozef.reisinger/.kube/config")
   --namespace value, -n value   kubernetes namespace (default: "default")
   --help, -h                    show help
```