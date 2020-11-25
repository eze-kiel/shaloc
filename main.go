package main

import "github.com/eze-kiel/shaloc/cmd"

var (
	// Version holds the build version
	Version string
	// BuildDate holds the build date
	BuildDate string
)

func main() {
	cmd.Execute()
}
