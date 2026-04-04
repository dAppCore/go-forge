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

	if err := run(*specPath, *outDir); err != nil {
		core.Print(os.Stderr, "forgegen: %v", err)
		os.Exit(1)
	}
}

func run(specPath, outDir string) error {
	spec, err := LoadSpec(specPath)
	if err != nil {
		return core.E("forgegen.main", "load spec", err)
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	core.Print(nil, "Loaded %d types, %d CRUD pairs", len(types), len(pairs))
	core.Print(nil, "Output dir: %s", outDir)

	if err := Generate(types, pairs, outDir); err != nil {
		return core.E("forgegen.main", "generate types", err)
	}
	return nil
}
