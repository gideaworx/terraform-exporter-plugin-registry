package main

import (
	"github.com/alecthomas/kong"
	"github.com/gideaworx/terraform-exporter-plugin-registry/backend/commands"
)

type options struct {
	CompileRegistry *commands.CompileRegistryCommand `cmd:"" help:"Merge all plugin yamls into a single registry yaml used by the CLI and site"`
	BuildSite       *commands.BuildSiteCommand       `cmd:"" help:"Output a production version of the site as a directory or zip file"`
}

func main() {
	var cli options

	kctx := kong.Parse(&cli)

	kctx.FatalIfErrorf(kctx.Run())
}
