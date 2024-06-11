package main

import (
	"fmt"
	"os"

	apackage "github.com/tklauser/buildinfo/apackage"
)

func main() {
	as := apackage.AStruct{}

	fmt.Fprintln(os.Stderr, as)
}
