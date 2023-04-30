package envgen

import (
	"go/ast"
	"strings"
)

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
