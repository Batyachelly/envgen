package envgen

import (
	"bytes"
	"fmt"
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
	Envs []Env
}

type EnvFiles map[string]EnvFile

// Save saves parsed environment variables to files using the DefaultMarshal.
func (efs EnvFiles) Save(outputDir string) error {
	for name, file := range efs {
		pathName := filepath.Join(outputDir, name)

		if err := os.WriteFile(pathName, DefaultMarshal(file.Envs), 0600); err != nil { //nolint:gofumpt,gomnd
			return fmt.Errorf("failed save result file %s: %w", name, err)
		}
	}

	return nil
}

// SaveWithCustomMarshal saves parsed environment variables to files using the received function.
func (efs EnvFiles) SaveWithCustomMarshal(outputDir string, marshal func([]Env) []byte) error {
	for name, file := range efs {
		pathName := filepath.Join(outputDir, name)

		if err := os.WriteFile(pathName, marshal(file.Envs), 0600); err != nil { //nolint:gofumpt,gomnd
			return fmt.Errorf("failed save result file %s: %w", name, err)
		}
	}

	return nil
}

// DefaultMarshal default marshal function.
func DefaultMarshal(envs []Env) []byte {
	var resultFile bytes.Buffer

	for _, env := range envs {
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
