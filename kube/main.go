package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"

	"github.com/agext/levenshtein"
	"github.com/b4b4r07/kubetools/kube/command"
)

var (
	subcommands = []string{
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
)

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	if args[0] == "kubectl" {
		return run("kubectl", args[1:])
	}

	if contains(subcommands, args[0]) {
		return run("kubectl", args)
	}

	if path, err := searchCommand(args[0]); err == nil {
		// Found user-defined command
		return runWithTTY(path, args[1:])
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
		fmt.Fprintf(os.Stderr, "%s: no such command\nThe most similar commands are %q", args[0], subs)
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

func searchCommand(cmd string) (path string, err error) {
	prefixes := []string{
		"kube", "kube-", "kubectl-",
	}
	for _, prefix := range prefixes {
		if path, err = exec.LookPath(prefix + cmd); err != nil {
			continue
		}
		return path, nil
	}
	return "", fmt.Errorf("%s: no such command", cmd)
}

func similarCommands(arg string) (found []string) {
	var max float64
	for _, sub := range subcommands {
		score := round(levenshtein.Similarity(sub, arg, nil) * 100)
		if score >= max {
			max = score
			if score > 65 {
				found = append(found, sub)
			}
		}
	}
	return found
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
