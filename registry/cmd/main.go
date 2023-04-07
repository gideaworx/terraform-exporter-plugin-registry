package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gideaworx/terraform-exporter-plugin-registry/registry"
	"gopkg.in/yaml.v3"
)

func main() {
	log.SetFlags(log.Default().Flags() | log.Lshortfile)
	var pluginRegistry registry.PluginRegistry

	thisDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	indexFile, err := os.Open(filepath.Join(thisDir, "index.yaml"))
	if err != nil {
		log.Fatal(err)
	}
	defer indexFile.Close()

	if err = yaml.NewDecoder(indexFile).Decode(&pluginRegistry); err != nil {
		log.Fatal(err)
	}

	pluginRegistry.Plugins = []registry.Plugin{}

	ymlFiles, err := filepath.Glob(filepath.Join(thisDir, "plugins", "*.yml"))
	if err != nil {
		log.Fatal(err)
	}

	yamlFiles, err := filepath.Glob(filepath.Join(thisDir, "plugins", "*.yaml"))
	if err != nil {
		log.Fatal(err)
	}

	pluginFiles := append(ymlFiles, yamlFiles...)
	for _, file := range pluginFiles {
		var subRegistry registry.PluginRegistry
		pluginFile, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer pluginFile.Close()

		if err = yaml.NewDecoder(pluginFile).Decode(&subRegistry); err != nil {
			log.Fatal(err)
		}

		pluginRegistry.Plugins = append(pluginRegistry.Plugins, subRegistry.Plugins...)
	}

	buildPath := os.Getenv("BUILD_PATH")
	if strings.TrimSpace(buildPath) == "" {
		buildPath = filepath.Join(thisDir, "..", "build", "web")
	}

	outFile, err := os.Create(filepath.Join(buildPath, "index.yaml"))
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if err = yaml.NewEncoder(outFile).Encode(pluginRegistry); err != nil {
		log.Fatal(err)
	}
}
