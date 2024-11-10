package presentation

import (
	"github.com/KarnerTh/mermerd/config"
	"github.com/fatih/color"
)

func ShowInfo(c config.MermerdConfig, value string) {
	if c.OutputMode() == config.Stdout {
		return
	}

	color.Blue(value)
}
