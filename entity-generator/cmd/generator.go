package main

import (
	"os"

	ts "github.com/mabels/wueste/entity-generator/ts"
)

func main() {
	ts.MainAction(os.Args[1:])
}
