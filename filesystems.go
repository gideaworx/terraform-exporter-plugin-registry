package terraformexporterpluginregistry

import "embed"

//go:embed registry/index.yaml registry/plugins
var PluginRegistries embed.FS

//go:embed frontend/.yarnrc.yml frontend/.yarn/releases frontend/.yarn/plugins frontend/public frontend/src frontend/yarn.lock frontend/package.json frontend/tsconfig.json
var RawSite embed.FS

//go:embed image/Dockerfile image/main.go
var ImageServer embed.FS
