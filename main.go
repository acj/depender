package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	outputPath := flag.String("output", "deps.go", "path to output source file")
	excludedPaths := flag.String("exclude", "", "package paths (or sub-strings) to exclude")
	flag.Parse()

	defaultExcludedPathSegments := []string{"internal"}
	userExcludedPathSegments := []string{}
	if *excludedPaths != "" {
		userExcludedPathSegments = strings.Split(*excludedPaths, ",")
	}
	excludedPathSegments := append(defaultExcludedPathSegments, userExcludedPathSegments...)

	m := Manifest{
		PackagePaths:         strings.Join(flag.Args(), " "),
		ExcludedPathSegments: excludedPathSegments,
	}
	buf, err := m.Generate()

	f, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("couldn't create output file '%s': %v", *outputPath, err)
		os.Exit(-1)
	}
	defer f.Close()

	if 	_, err := f.Write(buf); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(-1)
	}
}
