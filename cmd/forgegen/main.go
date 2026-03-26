package main

import (
	"flag"

	core "dappco.re/go/core"
)

func main() {
	specPath := flag.String("spec", "testdata/swagger.v1.json", "path to swagger.v1.json")
	outDir := flag.String("out", "types", "output directory for generated types")
	flag.Parse()

	spec, err := LoadSpec(*specPath)
	if err != nil {
		panic(core.E("forgegen.main", "load spec", err))
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	core.Print(nil, "Loaded %d types, %d CRUD pairs", len(types), len(pairs))
	core.Print(nil, "Output dir: %s", *outDir)

	if err := Generate(types, pairs, *outDir); err != nil {
		panic(core.E("forgegen.main", "generate types", err))
	}
}
