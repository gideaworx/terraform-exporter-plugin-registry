package registry

type PluginType string

const (
	Native PluginType = "native"
	Python PluginType = "python"
	NodeJS PluginType = "nodejs"
	Java   PluginType = "java"
)

type TargetArchitecture string

const (
	DarwinAmd64  TargetArchitecture = "darwin/amd64"
	DarwinArm64  TargetArchitecture = "darwin/arm64"
	LinuxAmd64   TargetArchitecture = "linux/amd64"
	LinuxArm64   TargetArchitecture = "linux/arm64"
	WindowsAmd64 TargetArchitecture = "windows/amd64"
	WindowsArm64 TargetArchitecture = "windows/arm64"
)

type pluginAuthor struct {
	Name    string `yaml:"name"`
	Email   string `yaml:"email"`
	Company string `yaml:"company,omitempty"`
}

type pluginExecutable struct {
	Locator string     `yaml:"locator"`
	Type    PluginType `yaml:"type"`
	Info    struct {
		ExtraArgs []string `yaml:"args,omitempty"`
		Checksum  string   `yaml:"sha256sum,omitempty"`
	} `yaml:"info"`
}

type pluginVersion struct {
	Version      string                                  `yaml:"version"`
	DownloadInfo map[TargetArchitecture]pluginExecutable `yaml:"download"`
}

type Plugin struct {
	Name        string          `yaml:"name"`
	Description string          `yaml:"description"`
	Homepage    string          `yaml:"homepage"`
	Authors     []pluginAuthor  `yaml:"authors"`
	Versions    []pluginVersion `yaml:"versions"`
}

type PluginRegistry struct {
	Name    string   `yaml:"name"`
	BaseURL string   `yaml:"baseURL"`
	Plugins []Plugin `yaml:"plugins,omitempty"`
}
