package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/Batyachelly/envgen"
)

func main() {
	target := flag.String("target", "", "target file for parsing")
	structs := flag.String("structs", "", "list of structures for parsing, separated by commas")
	outPutDir := flag.String("output_dir", "", "directory where the .env files will be generated")

	flag.Parse()

	if *target == "" {
		log.Fatal("error: empty target flag")
	}

	if *structs == "" {
		log.Fatal("error: empty structs flag")
	}

	if *outPutDir == "" {
		log.Fatal("error: empty outPutDir flag")
	}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, *target, nil, parser.ParseComments|parser.Trace)
	if err != nil {
		fmt.Println(err)
		return
	}

	structsNames := strings.Split(*structs, ",")
	if len(structsNames) == 0 {
		log.Fatal("invalid structs value")
	}

	ep, err := envgen.NewParser(f.Scope.Objects).FindStructs(structsNames)
	if err != nil {
		log.Fatal(err)
	}

	if err := ep.ParseFields().Save(*outPutDir); err != nil {
		log.Fatal(err)
	}
}
