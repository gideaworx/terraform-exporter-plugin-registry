package commands

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
)

const output = `# Copy this output and run it, or re-run the command wrapped in an "eval" statement
# You must have 'buildctl' and 'docker' installed for this to work.
buildctl build --frontend=dockerfile.v0 --local context=%[1]q --local dockerfile=%[1]q --output type=docker,name=%[4]q,dest=%[1]q/image.tar
[[ "%[2]t" == "true" ]] && docker load < %[1]q/image.tar
[[ "%[3]t" == "true" ]] && gzip < %[1]q/image.tar > ./image.tgz
rm %[1]q/image.tar
`

type BuildImageCommand struct {
	ImageName string `short:"n" default:"registry-server" help:"the name of the docker image in the daemon. ignored if --load-image is not set"`
	LoadImage bool   `help:"Load the image into the local docker daemon"`
	SaveImage bool   `help:"Save the image as a gzipped tarball on the filesystem"`
}

func (b *BuildImageCommand) Run(ctx *kong.Context, site embed.FS) error {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}

	p := &PackageSiteCommand{
		OutputDirectory: tempDir,
		Zip:             false,
	}

	if err := p.Run(site); err != nil {
		return err
	}

	mainGo, err := site.ReadFile("image-server/main.go")
	if err != nil {
		return err
	}

	dockerfile, err := site.ReadFile("image-server/Dockerfile")
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(tempDir, "main.go"), mainGo, 0666); err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module image-server\ngo 1.20"), 0666); err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(tempDir, "Dockerfile"), dockerfile, 0666); err != nil {
		return err
	}

	b.writeCommand(ctx.Stdout, tempDir)

	return nil
}

func (b *BuildImageCommand) writeCommand(out io.Writer, tempDir string) {
	fmt.Fprintf(out, output, tempDir, b.LoadImage, b.SaveImage, b.ImageName)
}
