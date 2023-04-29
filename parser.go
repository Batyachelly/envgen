package envgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
)

type env struct {
	name         string
	description  []string
	mandatory    bool
	defaultValue string
	example      string
	separator    string
}

type envFile struct {
	envs []env
}

func (ef envFile) marshal() []byte {
	b := bytes.Buffer{}
	for _, e := range ef.envs {
		for _, line := range e.description {
			b.WriteString("#" + line + "\n")
		}

		var additional []string

		if e.example != "" {
			additional = append(additional, "Пример:\""+e.example+"\"")
		}

		if e.separator != "" {
			additional = append(additional, "Разделитель:\""+e.separator+"\"")
		}

		if len(additional) > 0 {
			b.WriteString("#" + strings.Join(additional, ". ") + "\n")
		}

		if e.mandatory {
			b.WriteString("#")
		}

		b.WriteString(e.name + "=" + e.defaultValue + "\n\n")

	}

	return b.Bytes()
}

type envFiles map[string]envFile

func (efs envFiles) Save(outputDir string) error {
	for name, file := range efs {
		if err := os.WriteFile(filepath.Join(outputDir, name), file.marshal(), 0644); err != nil {
			return err
		}
	}

	return nil
}

type envParser struct {
	objects map[string]*ast.Object
}

func NewParser(objects map[string]*ast.Object) envParser {
	return envParser{
		objects: objects,
	}
}

func (ep envParser) FindStructs(targetStructs []string) (envParser, error) {
	objects := make(map[string]*ast.Object, len(targetStructs))

	for _, targetStruct := range targetStructs {
		for key, obj := range ep.objects {
			if obj.Name == targetStruct {
				if obj.Kind != ast.Typ {
					return ep, fmt.Errorf("invalid type of target: %s", targetStruct)
				}

				objects[key] = obj
			}
		}
	}

	return envParser{
		objects: objects,
	}, nil
}

func (ep envParser) ParseFields() envFiles {
	res := make(envFiles, len(ep.objects))

	for _, object := range ep.objects {
		if object.Decl == nil || object.Decl.(*ast.TypeSpec).Type == nil {
			continue
		}

		comments := object.Decl.(*ast.TypeSpec).Comment
		if comments == nil {
			continue
		}

		commentsList := comments.List

		if len(commentsList) == 0 {
			continue
		}

		fileName := strings.TrimSpace(strings.TrimPrefix(commentsList[0].Text, "//"))

		fields := object.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields
		if fields == nil {
			continue
		}

		fieldsList := fields.List

		file := res[fileName]

		if file.envs == nil {
			file.envs = make([]env, 0, len(fieldsList))
		}

		file.envs = append(file.envs, parseFields(fieldsList)...)

		res[fileName] = file
	}

	return res
}

func parseFields(fieldsList []*ast.Field) (envs []env) {
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

		e, ok := parseFromTagsString(field.Tag.Value)
		if !ok {
			continue
		}

		if field.Doc != nil {
			docList := field.Doc.List

			for _, doc := range docList {
				line := strings.TrimSpace(strings.TrimPrefix(doc.Text, "//"))
				e.description = append(e.description, line)
			}
		}

		envs = append(envs, e)
	}

	return envs
}

func parseFromTagsString(tagsString string) (e env, ok bool) {
	tagStrings := strings.Fields(tagsString[1 : len(tagsString)-1])

	for _, tagString := range tagStrings {
		delimiter := strings.Index(tagString, ":")
		if delimiter == -1 {
			return e, false
		}

		tagKey := tagString[:delimiter]
		tagValues := strings.Split(tagString[delimiter+2:len(tagString)-1], ",")

		if len(tagValues) == 0 {
			return e, false
		}

		switch tagKey {
		case "env":
			for _, value := range tagValues {
				switch value {
				case "required":
					e.mandatory = true
				default:
					e.name = value
				}
			}
		case "envDefault":
			e.defaultValue = tagValues[0]
		case "envSeparator":
			e.separator = tagValues[0]
		case "envExample":
			e.separator = tagValues[0]
		}
	}

	return e, true
}
