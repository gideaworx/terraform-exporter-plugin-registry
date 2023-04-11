package registry

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

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
	MultiArch    TargetArchitecture = "multi-arch"
)

type PluginAuthor struct {
	Name    string `yaml:"name"`
	Email   string `yaml:"email"`
	Company string `yaml:"company,omitempty"`
}

type PluginExecutable struct {
	Locator string     `yaml:"locator"`
	Type    PluginType `yaml:"type"`
	Info    struct {
		ExtraArgs []string `yaml:"args,omitempty"`
		Checksum  string   `yaml:"sha256sum,omitempty"`
	} `yaml:"info"`
}

type PluginVersion struct {
	Version      string                                  `yaml:"version"`
	DownloadInfo map[TargetArchitecture]PluginExecutable `yaml:"download"`
}

type Plugin struct {
	Name        string          `yaml:"name"`
	Description string          `yaml:"description"`
	Homepage    string          `yaml:"homepage"`
	LastUpdated ISO8601Time     `yaml:"lastUpdated"`
	Authors     []PluginAuthor  `yaml:"authors"`
	Versions    []PluginVersion `yaml:"versions"`
}

type PluginRegistry struct {
	Name    string   `yaml:"name"`
	BaseURL string   `yaml:"baseURL"`
	Plugins []Plugin `yaml:"plugins,omitempty"`
}

type ISO8601Time time.Time

func (i *ISO8601Time) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return nil
	}

	var nodeVal string
	switch node.Kind {
	case yaml.AliasNode:
		nodeVal = node.Alias.Value
	case yaml.ScalarNode:
		nodeVal = node.Value
	default:
		return fmt.Errorf("did not expect yaml node type %v", node.Kind)
	}

	t, err := time.Parse(time.RFC3339, nodeVal)
	if err != nil {
		return err
	}

	*i = ISO8601Time(t.UTC())
	return nil
}

func (i ISO8601Time) MarshalYAML() (interface{}, error) {
	t := time.Time(i)
	return t.UTC().Format(time.RFC3339), nil
}
