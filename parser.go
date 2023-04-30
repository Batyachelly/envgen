package envgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
)

// EnvParser implement envs parsing from .go file.
type EnvParser struct {
	objects map[string]*ast.Object
}

// NewParserFromASTObjects creates a new parser based on the AST object map.
func NewParserFromASTObjects(objects map[string]*ast.Object) EnvParser {
	return EnvParser{
		objects: objects,
	}
}

// NewParserFromFile creates a new parser based on the .go file.
func NewParserFromFile(filePath string) (EnvParser, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments|parser.Trace)
	if err != nil {
		return EnvParser{}, fmt.Errorf("failed to parse source: %w", err)
	}

	return NewParserFromASTObjects(f.Scope.Objects), nil
}

// NewParserFromReader creates a new parser based on data from the io.Reader.
func NewParserFromReader(r io.Reader) (EnvParser, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", r, parser.ParseComments|parser.Trace)
	if err != nil {
		return EnvParser{}, fmt.Errorf("failed to parse source: %w", err)
	}

	return NewParserFromASTObjects(f.Scope.Objects), nil
}

// FindStructs filters structures leaving only those specified in the argument.
func (ep EnvParser) FindStructs(targetStructs []string) (EnvParser, error) {
	objects := make(map[string]*ast.Object, len(targetStructs))

	for _, targetStruct := range targetStructs {
		for key, obj := range ep.objects {
			if obj.Name != targetStruct {
				continue
			}

			if obj.Kind != ast.Typ {
				return ep, fmt.Errorf("fiailde to parse struct %s: %w", targetStruct, ErrInvalidASTObjectType)
			}

			objects[key] = obj
		}
	}

	return EnvParser{
		objects: objects,
	}, nil
}

// ParseFields parse structs to EnvFiles.
func (ep EnvParser) ParseFields() EnvFiles {
	envFiles := make(EnvFiles, len(ep.objects))

	for _, object := range ep.objects {
		typeSpec, ok := object.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		if typeSpec.Comment == nil || len(typeSpec.Comment.List) == 0 {
			continue
		}

		comment := typeSpec.Comment.List[0].Text

		fileName := strings.TrimSpace(strings.TrimPrefix(comment, "//"))

		file := envFiles[fileName]

		if structType.Fields == nil {
			continue
		}

		fieldsList := structType.Fields.List

		if file.Envs == nil {
			file.Envs = make([]Env, 0, len(fieldsList))
		}

		file.Envs = append(file.Envs, parseFields(fieldsList)...)

		envFiles[fileName] = file
	}

	return envFiles
}
