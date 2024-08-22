package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func getKubeNamespaces() ([]string, error) {
	cmd := exec.Command("kubectl", "get", "namespaces", "-o", "custom-columns=:metadata.name")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	namespaces := strings.Split(strings.TrimSpace(string(out)), "\n")
	return namespaces, nil
}

func switchNamespace(namespace string) error {
	cmd := exec.Command("kubectl", "config", "set-context", "--current", "--namespace="+namespace)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	namespaces, err := getKubeNamespaces()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching namespaces: %v\n", err)
		os.Exit(1)
	}

	prompt := promptui.Select{
		Label: "Select namespace",
		Items: namespaces,
	}

	_, selectedNamespace, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting namespace: %v\n", err)
		os.Exit(1)
	}

	err = switchNamespace(selectedNamespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error switching namespace: %v\n", err)
	}

	fmt.Printf("Set current namespace to: %s\n", selectedNamespace)
}
