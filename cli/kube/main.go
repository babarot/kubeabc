package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/agext/levenshtein"
	"github.com/b4b4r07/kubeabc/cli/kube/command"
)

var subcommands = []string{
	"create",
	"expose",
	"run",
	"run-container",
	"set",
	"get",
	"explain",
	"edit",
	"delete",
	"rollout",
	"rolling-update",
	"rollingupdate",
	"scale",
	"resize",
	"autoscale",
	"certificate",
	"cluster-info",
	"clusterinfo",
	"top",
	"cordon",
	"uncordon",
	"drain",
	"taint",
	"describe",
	"logs",
	"attach",
	"exec",
	"port-forward",
	"proxy",
	"cp",
	"auth",
	"apply",
	"patch",
	"replace",
	"update",
	"convert",
	"label",
	"annotate",
	"completion",
	"api-versions",
	"config",
	"help",
	"plugin",
	"version",
}

var resources = []string{
	"all",
	"certificatesigningrequests", "certificatesigningrequest", "csr",
	"clusterrolebindings", "clusterrolebinding",
	"clusterroles", "clusterroles",
	"clusters", "cluster",
	"componentstatuses", "componentstatus", "cs",
	"configmaps", "configmap", "cm",
	"controllerrevisions", "controllerrevision",
	"cronjobs", "cronjob",
	"daemonsets", "daemonset", "ds",
	"deployments", "deployment", "deploy",
	"endpoints", "endpoint", "ep",
	"events", "event", "ev",
	"horizontalpodautoscalers", "horizontalpodautoscalers", "hpa",
	"ingresses", "ingress", "ing",
	"jobs", "job",
	"limitranges", "limitrange", "limits",
	"namespaces", "namespace", "ns",
	"networkpolicies", "networkpolicy", "netpol",
	"nodes", "node", "no",
	"persistentvolumeclaims", "persistentvolumeclaim", "pvc",
	"persistentvolumes", "persistentvolume", "pv",
	"poddisruptionbudgets", "poddisruptionbudget", "pdb",
	"podpreset",
	"pods", "pod", "po",
	"podsecuritypolicies", "podsecuritypolicy", "psp",
	"podtemplates", "podtemplate",
	"replicasets", "replicaset", "rs",
	"replicationcontrollers", "replicationcontroller", "rc",
	"resourcequotas", "resourcequotas", "quota",
	"rolebindings", "rolebinding",
	"roles", "role",
	"secrets", "secret",
	"serviceaccounts", "serviceaccount", "sa",
	"services", "service", "svc",
	"statefulsets", "statefulset",
	"storageclasses", "storageclass",
	"thirdpartyresources", "thirdpartyresource",
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	if len(args) == 0 {
		return run("kubectl", []string{"help"})
	}

	if len(args) > 1 {
		results := similarResources(args[1])
		if !contains(results, args[1]) {
			switch len(results) {
			case 0:
				// through
			case 1:
				fmt.Fprintf(os.Stdout,
					"You called a k8s resource named '%s', which does not exist.\nContinuing under the assumption that you meant '%s'\n",
					args[1], results[0])
				args[1] = results[0]
			default:
				fmt.Fprintf(os.Stderr,
					"%s: no such resource\nThe most similar resources are %q\n",
					args[0], results)
			}
		}
	}

	if args[0] == "kubectl" {
		return run("kubectl", args[1:])
	}

	if contains(subcommands, args[0]) {
		return run("kubectl", args)
	}

	cmds := searchCommands(args[0])
	switch len(cmds) {
	case 0:
		// through
	case 1:
		return runWithTTY(cmds[0], args[1:])
	default:
		fmt.Fprintf(os.Stderr,
			"Some commands are found: %q\n",
			cmds)
		return 1
	}

	subs := similarCommands(args[0])
	switch len(subs) {
	case 0:
		// through
	case 1:
		fmt.Fprintf(os.Stdout,
			"You called a kubectl command named '%s', which does not exist.\nContinuing under the assumption that you meant '%s'\n",
			args[0], subs[0])
		args[0] = subs[0]
		return run("kubectl", args)
	default:
		fmt.Fprintf(os.Stderr,
			"%s: no such command\nThe most similar commands are %q\n",
			args[0], subs)
		return 1
	}

	fmt.Fprintf(os.Stderr, "%s: no such command in kubectl\n", args[0])
	return 1
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func run(arg string, args []string) int {
	switch args[0] {
	case "apply", "delete":
		prompt()
	default:
	}
	c := command.New(command.Join(arg, args))
	if err := c.Run(); err != nil {
		// Unexpected error
		log.Fatal(err)
	}
	res := c.Result()
	if res.Failed {
		fmt.Fprintf(os.Stderr, "Error: %v\n", res.StderrString())
		return res.ExitCode
	}
	out := res.StdoutString()
	if len(out) > 0 {
		fmt.Fprintln(os.Stdout, out)
	}
	return res.ExitCode
}

func runWithTTY(arg string, args []string) int {
	c := command.New(command.Join(arg, args))
	if err := c.RunWithTTY(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		return 1
	}
	return 0
}

func searchCommands(arg string) (results []string) {
	prefixes := []string{
		"kube", "kube-", "kubectl-",
	}
	for _, prefix := range prefixes {
		cmd := prefix + arg
		if _, err := exec.LookPath(cmd); err != nil {
			continue
		}
		results = append(results, cmd)
	}
	return
}

func similarCommands(arg string) (results []string) {
	var max float64
	for _, cmd := range subcommands {
		score := round(levenshtein.Similarity(cmd, arg, nil) * 100)
		if score >= max {
			max = score
			if score > 65 {
				results = append(results, cmd)
			}
		}
	}
	return
}

func similarResources(arg string) (results []string) {
	var max float64
	for _, resource := range resources {
		score := round(levenshtein.Similarity(resource, arg, nil) * 100)
		if score >= max {
			max = score
			if score > 65 {
				results = append(results, resource)
			}
		}
	}
	return
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}

func prompt() {
	file := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.LoadFromFile(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Press Return key to continue\n-> current context %q", config.CurrentContext)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
