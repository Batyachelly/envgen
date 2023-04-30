package envgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Env struct {
	name         string
	description  []string
	mandatory    bool
	defaultValue string
	example      string
	separator    string
}

type EnvFile struct {
	envs []Env
}

func (ef EnvFile) marshal() []byte {
	var resultFile bytes.Buffer

	for _, env := range ef.envs {
		for _, line := range env.description {
			resultFile.WriteString("#" + line + "\n")
		}

		var additional []string

		if env.example != "" {
			additional = append(additional, "Пример:\""+env.example+"\"")
		}

		if env.separator != "" {
			additional = append(additional, "Разделитель:\""+env.separator+"\"")
		}

		if len(additional) > 0 {
			resultFile.WriteString("#" + strings.Join(additional, ". ") + "\n")
		}

		if env.mandatory {
			resultFile.WriteString("#")
		}

		resultFile.WriteString(env.name + "=" + env.defaultValue + "\n\n")
	}

	return resultFile.Bytes()
}

type EnvFiles map[string]EnvFile

func (efs EnvFiles) Save(outputDir string) error {
	for name, file := range efs {
		pathName := filepath.Join(outputDir, name)

		if err := os.WriteFile(pathName, file.marshal(), 0600); err != nil { //nolint:gofumpt,gomnd
			return fmt.Errorf("failed save result file %s: %w", name, err)
		}
	}

	return nil
}

type EnvParser struct {
	objects map[string]*ast.Object
}

func NewParserFromASTObjects(objects map[string]*ast.Object) EnvParser {
	return EnvParser{
		objects: objects,
	}
}

func NewParserFromFile(filePath string) (EnvParser, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments|parser.Trace)
	if err != nil {
		return EnvParser{}, fmt.Errorf("failed to parse source: %w", err)
	}

	return NewParserFromASTObjects(f.Scope.Objects), nil
}

func NewParserFromReader(r io.Reader) (EnvParser, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", r, parser.ParseComments|parser.Trace)
	if err != nil {
		return EnvParser{}, fmt.Errorf("failed to parse source: %w", err)
	}

	return NewParserFromASTObjects(f.Scope.Objects), nil
}

func (ep EnvParser) FindStructs(targetStructs []string) (EnvParser, error) {
	objects := make(map[string]*ast.Object, len(targetStructs))

	for _, targetStruct := range targetStructs {
		for key, obj := range ep.objects {
			if obj.Name == targetStruct {
				if obj.Kind != ast.Typ {
					return ep, fmt.Errorf("fiailde to parse struct %s: %w", targetStruct, ErrInvalidType)
				}

				objects[key] = obj
			}
		}
	}

	return EnvParser{
		objects: objects,
	}, nil
}

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

		comments := typeSpec.Comment
		if comments == nil {
			continue
		}

		commentsList := comments.List

		if len(commentsList) == 0 {
			continue
		}

		fileName := strings.TrimSpace(strings.TrimPrefix(commentsList[0].Text, "//"))

		file := envFiles[fileName]

		fields := structType.Fields
		if fields == nil {
			continue
		}

		fieldsList := fields.List

		if file.envs == nil {
			file.envs = make([]Env, 0, len(fieldsList))
		}

		file.envs = append(file.envs, parseFields(fieldsList)...)

		envFiles[fileName] = file
	}

	return envFiles
}

func parseFields(fieldsList []*ast.Field) []Env {
	envs := make([]Env, 0, len(fieldsList))

	for _, field := range fieldsList {
		switch v := field.Type.(type) {
		case *ast.Ident:
		case *ast.SelectorExpr:
		case *ast.ArrayType:
		case *ast.StructType:
			envs = append(envs, parseFields(v.Fields.List)...)

			continue
		}

		if field.Tag == nil {
			continue
		}

		env, ok := parseEnvFromTagsString(field.Tag.Value)
		if !ok {
			continue
		}

		if field.Doc != nil {
			docList := field.Doc.List

			for _, doc := range docList {
				line := strings.TrimSpace(strings.TrimPrefix(doc.Text, "//"))
				env.description = append(env.description, line)
			}
		}

		envs = append(envs, env)
	}

	return envs
}

func parseEnvFromTagsString(tagsString string) (Env, bool) {
	var env Env

	tagStrings := strings.Fields(tagsString[1 : len(tagsString)-1])

	for _, tagString := range tagStrings {
		delimiter := strings.Index(tagString, ":")
		if delimiter == -1 {
			return env, false
		}

		tagKey := tagString[:delimiter]
		tagValues := strings.Split(tagString[delimiter+2:len(tagString)-1], ",")

		if len(tagValues) == 0 {
			return env, false
		}

		switch tagKey {
		case "env":
			env.name = tagValues[0]
			env.mandatory = contains(tagValues, "required")
		case "envDefault":
			env.defaultValue = tagValues[0]
		case "envSeparator":
			env.separator = tagValues[0]
		case "envExample":
			env.separator = tagValues[0]
		}
	}

	return env, true
}
