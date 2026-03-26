package main

import (
	"flag"
	"os"

	core "dappco.re/go/core"
)

func main() {
	specPath := flag.String("spec", "testdata/swagger.v1.json", "path to swagger.v1.json")
	outDir := flag.String("out", "types", "output directory for generated types")
	flag.Parse()

	spec, err := LoadSpec(*specPath)
	if err != nil {
		core.Print(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	core.Print(os.Stdout, "Loaded %d types, %d CRUD pairs\n", len(types), len(pairs))
	core.Print(os.Stdout, "Output dir: %s\n", *outDir)

	if err := Generate(types, pairs, *outDir); err != nil {
		core.Print(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
