package util

import "github.com/fatih/color"

func ShowSuccess() {
	color.Green(`

Diagram was created successfully

`)
}

func ShowError() {
	color.Red(`

X Something went wrong!

`)
}
