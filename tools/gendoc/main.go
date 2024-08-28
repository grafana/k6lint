// Package main contains CLI documentation generator tool.
package main

import (
	"github.com/grafana/clireadme"
	"github.com/grafana/k6lint/cmd"
)

func main() {
	root, _ := cmd.New()
	clireadme.Main(root, 1)
}
