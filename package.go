package main

import (
	"golang.org/x/tools/go/packages"
	"sort"
	"strings"
)

func dependenciesForPackages(paths string, excludedPathSegments []string) ([]string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps,
	}
	loadedPackages, err := packages.Load(cfg, paths)
	if err != nil {
		return nil, err
	}

	dependencies := []string{}
	for _, p := range loadedPackages {
		for _, i := range p.Imports {
			if isAllowedPackagePath(i.PkgPath, excludedPathSegments) && isNonStandardDependency(i.PkgPath) {
				dependencies = append(dependencies, i.PkgPath)
			}
		}
	}

	dependencies = unique(dependencies)
	sort.Strings(dependencies)

	return dependencies, nil
}

func isAllowedPackagePath(packagePath string, excludedPathSegments []string) bool {
	for _, p := range excludedPathSegments {
		if strings.Contains(packagePath, p) {
			return false
		}
	}
	return true
}

func isNonStandardDependency(packagePath string) bool {
	return !standardLibraryPackages[packagePath]
}

func unique(strs []string) []string {
	seenStrings := make(map[string]struct{})
	uniqueStrings := []string{}
	for _, s := range strs {
		if _, ok := seenStrings[s]; !ok {
			seenStrings[s] = struct{}{}
			uniqueStrings = append(uniqueStrings, s)
		}
	}
	return uniqueStrings
}