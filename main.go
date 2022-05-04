// Playing around with https://pkg.go.dev/debug/buildinfo

package main

import (
	"debug/buildinfo"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 0 {
		fmt.Fprintf(os.Stderr, "usage: buildinfo binary...\n")
		os.Exit(1)
	}

	for i, f := range os.Args[1:] {
		info, err := buildinfo.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read build info from %s: %s\n", f, err)
			continue
		}

		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("Go binary: %s\n", f)
		fmt.Printf("Compiled with Go version: %s\n", info.GoVersion)
		fmt.Printf("Build settings:\n")
		for _, s := range info.Settings {
			fmt.Printf("  %s: %s\n", s.Key, s.Value)
		}
		for _, d := range info.Deps {
			fmt.Printf("  %s@%s %s\n", d.Path, d.Version, d.Sum)
			if r := d.Replace; r != nil {
				fmt.Printf("    replaced by %s@%s %s\n", r.Path, r.Version, r.Sum)
			}
		}
	}
}
