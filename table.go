package main

import (
	"github.com/soerlemans/table/util"
)

func main() {
	err := initArgs()
	util.FailIf(err)
}
