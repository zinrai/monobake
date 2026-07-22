package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ExitOK          = 0
	ExitConfigError = 1
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage: monobake [options]

Resolve build target and version from Git tag.

Output:
  <target> <version>    Space-separated target name and version
  (empty)               No output if tag format is invalid

Options:`)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
Example:
  read TARGET VERSION <<< $(monobake -tag refs/tags/backend/v1.0.0)
  [ -n "$TARGET" ] && docker buildx bake \
    --set="${TARGET}.tags=${REGISTRY}/${TARGET}:${VERSION}" \
    "$TARGET" --push`)
	}

	tagStr := flag.String("tag", "", "Git tag to parse (required)")
	bakeFile := flag.String("file", DefaultBakeFile, "Path to Bake file")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersion {
		printVersion()
		return ExitOK
	}

	// No tag means nothing to build
	if *tagStr == "" {
		return ExitOK
	}

	// Parse the tag
	info, err := ParseGitRef(*tagStr)
	if err != nil {
		// Invalid tag format, nothing to build
		return ExitOK
	}

	// Load Bake file
	bake, err := LoadBakeFile(*bakeFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitConfigError
	}

	// Verify target exists
	if !bake.HasTarget(info.Target) {
		fmt.Fprintf(os.Stderr, "target %q not defined in %s\n", info.Target, *bakeFile)
		return ExitConfigError
	}

	fmt.Printf("%s %s\n", info.Target, info.Version)
	return ExitOK
}
