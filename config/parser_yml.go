package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/rainu/go-yacl"
)

var yamlLookupLocations = func() (result []string) {
	result = append(result, "/"+path.Join("etc", ".mcp-system-control.yml"))
	result = append(result, "/"+path.Join("etc", ".mcp-system-control.yaml"))
	result = append(result, "/"+path.Join("etc", "mcp-system-control", "config.yml"))
	result = append(result, "/"+path.Join("etc", "mcp-system-control", "config.yaml"))
	result = append(result, "/"+path.Join("usr", "local", "etc", ".mcp-system-control.yml"))
	result = append(result, "/"+path.Join("usr", "local", "etc", ".mcp-system-control.yaml"))
	result = append(result, "/"+path.Join("usr", "local", "etc", "mcp-system-control", "config.yml"))
	result = append(result, "/"+path.Join("usr", "local", "etc", "mcp-system-control", "config.yaml"))

	if home, err := os.UserHomeDir(); err == nil {
		result = append(result, path.Join(home, ".mcp-system-control.yml"))
		result = append(result, path.Join(home, ".mcp-system-control.yaml"))
		result = append(result, path.Join(home, ".config", ".mcp-system-control.yml"))
		result = append(result, path.Join(home, ".config", ".mcp-system-control.yaml"))
		result = append(result, path.Join(home, ".config", "mcp-system-control", "config.yml"))
		result = append(result, path.Join(home, ".config", "mcp-system-control", "config.yaml"))
	}

	binDir := path.Dir(os.Args[0])
	result = append(result, path.Join(binDir, ".mcp-system-control.yml"))
	result = append(result, path.Join(binDir, ".mcp-system-control.yaml"))

	if wd, err := os.Getwd(); err == nil {
		result = append(result, path.Join(wd, ".mcp-system-control.yml"))
		result = append(result, path.Join(wd, ".mcp-system-control.yaml"))
	}

	return
}

func processYamlFiles(config *yacl.Config, configFilePath string) {
	for _, location := range yamlLookupLocations() {
		processYamlFile(config, location)
	}
	if configFilePath != "" {
		processYamlFile(config, configFilePath)
	}
}

func processYamlFile(config *yacl.Config, path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	slog.Debug("Processing yaml file", "file", path)
	err = processYaml(config, f)
	if err != nil {
		panic(fmt.Errorf("unable to process yaml file %s: %w", path, err))
	}
}

func processYaml(config *yacl.Config, source io.Reader) error {
	err := config.ParseYaml(source)
	if err != nil {
		return fmt.Errorf("error while decoding yaml: %w", err)
	}
	return nil
}
