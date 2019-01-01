package main

import (
	"os"

	"github.com/fnuva/sqsbeat/cmd"

	_ "github.com/fnuva/sqsbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

