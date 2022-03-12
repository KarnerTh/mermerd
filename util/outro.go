package util

import (
	"fmt"
	"github.com/fatih/color"
	"mermerd/config"
)

func ShowSuccess() {
	color.Green(fmt.Sprintf(`

✓ Diagram was created successfully (%s)

`, config.OutputFileName()))
}

func ShowError() {
	color.Red(`

X Something went wrong!

`)
}
