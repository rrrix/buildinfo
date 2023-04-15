package main

import (
	"fmt"
	apackage "github.com/tklauser/buildinfo/apackage"
	"os"
)

func main() {

	as := apackage.AStruct{}

	fmt.Fprintln(os.Stderr, as)
}
