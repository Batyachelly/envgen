package main

import (
	"flag"
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

	structsNames := strings.Split(*structs, ",")
	if len(structsNames) == 0 {
		log.Fatal("invalid structs value")
	}

	envParser, err := envgen.NewParserFromFile(*target)
	if err != nil {
		log.Fatal(err)
	}

	envParser, err = envParser.FindStructs(structsNames)
	if err != nil {
		log.Fatal(err)
	}

	if err := envParser.ParseFields().Save(*outPutDir); err != nil {
		log.Fatal(err)
	}
}
