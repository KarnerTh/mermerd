package main

import (
	"flag"
	"fmt"
	"mermerd/analyzer"
	"mermerd/config"
	"mermerd/diagram"
	"mermerd/util"
)

func main() {
	err := config.LoadConfigFile()
	if err != nil {
		fmt.Println(err.Error())
		util.ShowError()
		return
	}

	flag.BoolVar(&config.ShowAllConstraints, "ac", false, "(allConstraints - default: false) Contain all constraints of selected tables, even though the table of the resulting constraint was not selected")
	flag.StringVar(&config.Schema, "s", "", "(schema - default: asks) The schema that should be used")
	flag.StringVar(&config.ConnectionString, "c", "", "(connectionString - default: asks) The connection string that should be used")
	flag.Parse()

	util.ShowIntro()
	result, err := analyzer.Analyze()
	if err != nil {
		fmt.Println(err.Error())
		util.ShowError()
		return
	}

	err = diagram.Create(result)
	if err != nil {
		fmt.Println(err.Error())
		util.ShowError()
		return
	}

	util.ShowSuccess()
}
