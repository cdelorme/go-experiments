package main

import (
	"os"
	"path"
)

// build variables can be separated from the func main file
var abc string
var softwareName string
var buildVersion string
var buildDate string

// an init function that can set "sane-defaults" where expected
func init() {
	if len(softwareName) == 0 {
		softwareName = path.Base(os.Args[0])
	}
}
