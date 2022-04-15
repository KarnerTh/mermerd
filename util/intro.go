package util

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func ShowIntro() {
	color.Green(fmt.Sprintf(`
.  ..___.__ .  ..___.__ .__ 
|\/|[__ [__)|\/|[__ [__)|  \
|  |[___|  \|  |[___|  \|__/ (%s)

Create MermaidJs diagrams from your database

`, viper.Get("version")))
}
