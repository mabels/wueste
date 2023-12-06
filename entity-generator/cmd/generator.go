package main

import (
	"os"

	ts "github.com/mabels/wueste/entity-generator/ts"
)

// GitCommit is injected during compile time
var GitCommit string

// Version is injected during compile time
var Version string

func main() {
	ts.MainAction(os.Args[1:], Version, GitCommit)
}
