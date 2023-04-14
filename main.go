// Playing around with https://pkg.go.dev/debug/buildinfo

package main

import (
	"debug/buildinfo"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
)

var BuildDate = "1970-01-01T00:00:00Z"
var GitCommit = "UnsetVar"
var GitSummary = "UnsetVar"
var GitBranch = "UnsetVar"

type BuildInfoFlags struct {
	self          *bool // display this binary's own build info
	showChecksums *bool // display checksums for each module
	raw           *bool // dump as debug.BuildInfo.String()
}

var buildInfoFlags = BuildInfoFlags{
	self:          flag.Bool("self", false, "display this binary's own build info"),
	showChecksums: flag.Bool("checksums", false, "display checksums for each module"),
	raw:           flag.Bool("raw", false, "display checksums for each module"),
}

func usage() {
	executableName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s [-h | --help] | [--checksums] [--raw] [[--self] | [binary [...]]]\n", executableName)
	fmt.Fprintf(os.Stderr, "  -h | --help: display this help message\n")
	fmt.Fprintf(os.Stderr, "  --checksums: display checksums for each module\n")
	fmt.Fprintf(os.Stderr, "  --raw: dump as debug.BuildInfo.String()\n")
	fmt.Fprintf(os.Stderr, "  --self: display this binary's own build info\n")
	fmt.Fprintf(os.Stderr, "  binary: display build info for each specified binary\n")
}

type Module struct {
	debug.Module
}

func (m Module) String() string {
	switch {
	case *buildInfoFlags.showChecksums == true && m.Replace == nil:
		return fmt.Sprintf("%s@%s %s", m.Path, m.Version, m.Sum)
	case *buildInfoFlags.showChecksums == true && m.Replace != nil:
		return fmt.Sprintf("%s@%s %s => %s@%s %s", m.Path, m.Version, m.Sum, m.Replace.Path, m.Replace.Version, m.Replace.Sum)
	case m.Replace != nil:
		return fmt.Sprintf("%s@%s => %s@%s", m.Path, m.Version, m.Replace.Path, m.Replace.Version)
	default:
		return fmt.Sprintf("%s@%s", m.Path, m.Version)
	}
}

func printBuildInfo(f string, info *debug.BuildInfo) {
	if *buildInfoFlags.raw {
		fmt.Printf("// %s\n", f)
		fmt.Println(info.String())
		return
	}
	fmt.Printf("Go binary: %s\n", f)
	fmt.Printf("Compiled with Go version: %s\n", info.GoVersion)
	fmt.Printf("Package: %s\n", info.Path)
	fmt.Printf("Main module: %s\n", Module{info.Main})
	fmt.Printf("Build settings:\n")
	for _, s := range info.Settings {
		fmt.Printf("  %s: %s\n", s.Key, s.Value)
	}
	fmt.Printf("Dependencies:\n")
	for _, d := range info.Deps {
		fmt.Printf("  %s\n", Module{*d})
	}
}

func printSelfBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintf(os.Stderr, "failed to read own embedded build info\n")
		os.Exit(1)
	}
	executablePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get absolute path of own executable: %v\n", err)
		executablePath = os.Args[0]
	}
	printBuildInfo(executablePath, info)
}

func main() {
	flag.CommandLine.Usage = usage
	flag.Parse()

	// os.Args[0] is the name of the executable
	if len(os.Args) == 1 || !*buildInfoFlags.self && len(flag.CommandLine.Args()) == 0 {
		usage()
		os.Exit(1)
	}

	if *buildInfoFlags.self && len(flag.CommandLine.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Error: cannot specify --self and a binary name\n")
		fmt.Fprintf(os.Stderr, "\n")
		usage()
		os.Exit(1)
	}

	if *buildInfoFlags.self {
		printSelfBuildInfo()
		return
	}

	for i, f := range flag.CommandLine.Args() {
		info, err := buildinfo.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read build info from %s: %v\n", f, err)
			continue
		}

		if i > 0 {
			fmt.Println()
		}
		printBuildInfo(f, info)
	}
}
