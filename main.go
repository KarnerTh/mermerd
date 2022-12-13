package main

/// TODO:  test
/// - tests
/// - TableNameResponse model?
/// - Aufteilung
/// - config abwärtskompatibel?

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/KarnerTh/mermerd/cmd"
)

// ldflags flags from goreleaser
var (
	version = "dev"
	commit  = "none"
)

func main() {
  test := "abc"
	viper.Set("version", version)
	viper.Set("commit", commit)
  fmt.Print(test)

	cmd.Execute()
}
