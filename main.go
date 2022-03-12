package main

import (
	"github.com/spf13/viper"
	"mermerd/cmd"
)

// ldflags flags from goreleaser
var (
	version = "dev"
	commit  = "none"
)

func main() {
	viper.Set("version", version)
	viper.Set("commit", commit)

	cmd.Execute()
}
