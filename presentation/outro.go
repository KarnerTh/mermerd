package presentation

import (
	"fmt"

	"github.com/KarnerTh/mermerd/config"
	"github.com/fatih/color"
)

func ShowSuccess(c config.MermerdConfig, fileName string) {
	if c.OutputMode() == config.Stdout {
		return
	}

	color.Green(fmt.Sprintf(`

✓ Diagram was created successfully (%s)

`, fileName))
}

func ShowError() {
	color.Red(`

X Something went wrong!

`)
}
