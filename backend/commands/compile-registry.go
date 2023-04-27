package commands

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

			lastUpdated := c.getLastModified(info.Name())
			for _, plugin := range pluginDef.Plugins {
				plugin.LastUpdated = lastUpdated
				pluginRegistry.Plugins = append(pluginRegistry.Plugins, plugin)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return yaml.NewEncoder(c.out).Encode(pluginRegistry)
}

func (c *CompileRegistryCommand) getLastModified(pluginFileName string) registry.ISO8601Time {
	git, err := exec.LookPath("git")
	if err != nil {
		log.Println("git not found on $PATH, returning current time")
		return registry.ISO8601Time(time.Now().UTC())
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := exec.Command(git, "rev-parse", "--show-toplevel")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = bytes.NewReader([]byte{})

	err = cmd.Run()
	if err != nil {
		log.Printf("git command failed, returning current time. error: %s", stderr.String())
		return registry.ISO8601Time(time.Now().UTC())
	}

	rootDir := strings.TrimSpace(stdout.String())
	pluginFile := filepath.Join(rootDir, "registry", "plugins", pluginFileName)

	info, err := os.Stat(pluginFile)
	if err != nil {
		log.Printf("could not stat %s, returning current time. error: %v", pluginFile, err)
		return registry.ISO8601Time(time.Now().UTC())
	}

	return registry.ISO8601Time(info.ModTime())
}
