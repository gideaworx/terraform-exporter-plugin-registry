package commands

import (
	"io"
	"io/fs"
	"log"
	"strings"

	"github.com/alecthomas/kong"
	terraformexporterpluginregistry "github.com/gideaworx/terraform-exporter-plugin-registry"
	"github.com/gideaworx/terraform-exporter-plugin-registry/backend/registry"
	"gopkg.in/yaml.v3"
)

type CompileRegistryCommand struct {
	out io.Writer `kong:"-"`
}

func (c *CompileRegistryCommand) BeforeApply(kctx *kong.Context) error {
	c.out = kctx.Stdout
	return nil
}

func (c *CompileRegistryCommand) Run() error {
	plugins, err := fs.Sub(terraformexporterpluginregistry.PluginRegistries, "registry")
	if err != nil {
		return err
	}

	var pluginRegistry registry.PluginRegistry
	indexFile, err := plugins.Open("index.yaml")
	if err != nil {
		return err
	}
	defer indexFile.Close()

	if err = yaml.NewDecoder(indexFile).Decode(&pluginRegistry); err != nil {
		return err
	}

	pluginRegistry.Plugins = []registry.Plugin{}
	if err = fs.WalkDir(plugins, "plugins", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() != "plugins" {
			return fs.SkipDir
		}

		lowerName := strings.ToLower(d.Name())
		if strings.HasSuffix(lowerName, ".yaml") || strings.HasSuffix(lowerName, ".yml") {
			info, err := d.Info()
			if err != nil {
				log.Println(err)
				return nil
			}

			var pluginDef registry.PluginRegistry
			pluginFile, err := plugins.Open(path)
			if err != nil {
				log.Println(err)
				return nil
			}
			defer pluginFile.Close()

			if err = yaml.NewDecoder(pluginFile).Decode(&pluginDef); err != nil {
				log.Println(err)
				return nil
			}

			lastUpdated := info.ModTime()
			for _, plugin := range pluginDef.Plugins {
				plugin.LastUpdated = registry.ISO8601Time(lastUpdated)
				pluginRegistry.Plugins = append(pluginRegistry.Plugins, plugin)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return yaml.NewEncoder(c.out).Encode(pluginRegistry)
}
