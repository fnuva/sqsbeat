package cmd

import (
	"github.com/fnuva/sqsbeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
)

// Name of this beat
var Name = "sqsbeat"

// RootCmd to handle fnuva cli
var RootCmd = cmd.GenRootCmd(Name, "", beater.New)
