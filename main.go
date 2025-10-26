package main

import (
	"fmt"
	"os"

	"github.com/reggles44/kubewatch/cmd"
)

var version = "0.0.4"

func main() {
	root, err := cmd.NewCmd()
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}

	if err := root.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}
