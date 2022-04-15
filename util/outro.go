package util

import (
	"fmt"

	"github.com/fatih/color"
)

func ShowSuccess(fileName string) {
	color.Green(fmt.Sprintf(`

✓ Diagram was created successfully (%s)

`, fileName))
}

func ShowError() {
	color.Red(`

X Something went wrong!

`)
}
