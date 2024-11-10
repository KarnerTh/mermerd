package presentation

import (
	"fmt"

	"github.com/KarnerTh/mermerd/config"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func ShowIntro(c config.MermerdConfig) {
	if c.OutputMode() == config.Stdout {
		return
	}

	color.Green(fmt.Sprintf(`
.  ..___.__ .  ..___.__ .__ 
|\/|[__ [__)|\/|[__ [__)|  \
|  |[___|  \|  |[___|  \|__/ (%s)

Create MermaidJs diagrams from your database

`, viper.Get("version")))
}
