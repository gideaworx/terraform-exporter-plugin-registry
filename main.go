package main

import (
	"embed"

	"github.com/alecthomas/kong"
	"github.com/gideaworx/terraform-exporter-plugin-registry/commands"
)

//go:embed build/web image-server/main.go image-server/Dockerfile
var site embed.FS

type options struct {
	Serve       *commands.ServeCommand       `cmd:"" help:"Host the registry server locally"`
	BuildImage  *commands.BuildImageCommand  `cmd:"" help:"Build an OCI image to run the site as a container"`
	PackageSite *commands.PackageSiteCommand `cmd:"" help:"Output the site as a directory or zip file"`
}

func main() {
	var cli options

	kctx := kong.Parse(&cli)
	kctx.Bind(site)

	kctx.FatalIfErrorf(kctx.Run())
}
