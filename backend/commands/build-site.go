package commands

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	terraformexporterpluginregistry "github.com/gideaworx/terraform-exporter-plugin-registry"
	"github.com/gideaworx/terraform-exporter-plugin-registry/backend/utils"
)

type BuildSiteCommand struct {
	OutputDirectory string `short:"o" help:"Where to write the site" default:"." type:"existingdir"`
	Keep            bool   `help:"keep temporary files"`
	Zip             bool   `help:"Zip the file before writing"`
	output          io.Writer
	zipfile         *zip.Writer
}

func (c *BuildSiteCommand) Run(ctx *kong.Context) error {
	c.output = ctx.Stderr
	buildTmp, err := os.MkdirTemp(os.TempDir(), "tfepr-fe-")
	if err != nil {
		return err
	}
	defer func() {
		if !c.Keep {
			os.RemoveAll(buildTmp)
		} else {
			fmt.Println(buildTmp)
		}
	}()

	if err = c.buildProductionSite(buildTmp); err != nil {
		return err
	}

	if c.Zip {
		zipFile, err := os.Create(filepath.Join(c.OutputDirectory, "registry-site.zip"))
		if err != nil {
			return err
		}
		defer zipFile.Close()

		c.zipfile = zip.NewWriter(zipFile)
		defer c.zipfile.Close()
	}

	builtDir := filepath.Join(buildTmp, "build")
	return filepath.WalkDir(builtDir, func(path string, d fs.DirEntry, e error) error {
		relative := strings.TrimPrefix(path, builtDir)

		if d.IsDir() {
			return c.mkdir(relative, 0777)
		}

		contents, err := os.ReadFile(filepath.Join(builtDir, relative))
		if err != nil {
			return err
		}

		return c.writeFile(relative, 0666, contents)
	})
}

func (c *BuildSiteCommand) buildProductionSite(buildTmp string) error {
	subFS, err := fs.Sub(terraformexporterpluginregistry.RawSite, "frontend")
	if err != nil {
		return err
	}

	if err = utils.ExportFS(subFS, buildTmp); err != nil {
		return err
	}

	node, err := exec.LookPath("node")
	if err != nil {
		return errors.New("could not find 'node' on PATH")
	}

	cmd := exec.Command(node, filepath.Join(buildTmp, ".yarn", "releases", "yarn-3.5.0.cjs"), "install", "--immutable")
	cmd.Stdout = c.output
	cmd.Stderr = c.output
	cmd.Dir = buildTmp
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command(node, filepath.Join(buildTmp, ".yarn", "releases", "yarn-3.5.0.cjs"), "build")
	cmd.Stdout = c.output
	cmd.Stderr = c.output
	cmd.Dir = buildTmp
	err = cmd.Run()
	if err != nil {
		return err
	}

	index, err := os.Create(filepath.Join(buildTmp, "build", "index.yaml"))
	if err != nil {
		return err
	}
	defer index.Close()

	compiler := &CompileRegistryCommand{
		out: index,
	}

	return compiler.Run()
}

func (c *BuildSiteCommand) mkdir(path string, mode os.FileMode) error {
	if c.Zip {
		_, err := c.zipfile.Create(filepath.Join(path, "/"))
		return err
	}

	return os.MkdirAll(filepath.Join(c.OutputDirectory, path), mode)
}

func (c *BuildSiteCommand) writeFile(path string, mode os.FileMode, contents []byte) error {
	var w io.Writer
	var err error
	if c.Zip {
		w, err = c.zipfile.Create(path)
	} else {
		w, err = os.OpenFile(filepath.Join(c.OutputDirectory, path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
		if err == nil {
			defer (w.(*os.File)).Close()
		}
	}

	if err != nil {
		return err
	}

	_, err = w.Write(contents)
	return err
}
