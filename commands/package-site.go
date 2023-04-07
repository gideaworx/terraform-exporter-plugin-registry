package commands

import (
	"archive/zip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type PackageSiteCommand struct {
	OutputDirectory string `short:"o" help:"Where to write the site" default:"." type:"existingdir"`
	Zip             bool   `help:"Zip the file before writing"`
	site            embed.FS
}

func (p *PackageSiteCommand) WriteZip() error {
	baseFile, err := os.Create(filepath.Join(p.OutputDirectory, "registry-site.zip"))
	if err != nil {
		return err
	}
	defer baseFile.Close()

	zipper := zip.NewWriter(baseFile)
	defer zipper.Close()

	return fs.WalkDir(p.site, ".", func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}

		path = strings.TrimPrefix(path, "build/web/")

		entryName := strings.TrimPrefix(path, "/")
		if d.IsDir() {
			entryName += "/"
			_, e := zipper.Create(entryName)
			return e
		}

		writer, e := zipper.Create(entryName)
		if e != nil {
			return e
		}

		f, e := p.site.Open(path)
		if e != nil {
			return e
		}
		defer f.Close()

		_, e = io.Copy(writer, f)
		return e
	})
}

func (p *PackageSiteCommand) WriteDir() error {
	baseDir := filepath.Join(p.OutputDirectory, "registry-site")
	if err := os.Mkdir(baseDir, 0o777); err != nil {
		return err
	}

	return fs.WalkDir(p.site, ".", func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}

		if strings.HasPrefix(path, "build/web/") {
			outpath := strings.TrimPrefix(path, "build/web/")

			if !d.IsDir() {
				if err := os.MkdirAll(filepath.Join(baseDir, filepath.Dir(outpath)), 0o777); err != nil {
					return err
				}

				info, err := d.Info()
				if err != nil {
					return err
				}

				outfile, err := os.OpenFile(filepath.Join(baseDir, outpath), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
				if err != nil {
					return err
				}
				defer outfile.Close()

				f, err := p.site.Open(path)
				if err != nil {
					return err
				}
				defer f.Close()

				written, err := io.Copy(outfile, f)
				if err != nil {
					return err
				}

				if written != info.Size() {
					return fmt.Errorf("%s is %d bytes but only %d bytes were written", path, info.Size(), written)
				}
			}
		}

		return nil
	})
}

func (p *PackageSiteCommand) Run(site embed.FS) error {
	p.site = site

	f := p.WriteDir
	if p.Zip {
		f = p.WriteZip
	}

	return f()
}
