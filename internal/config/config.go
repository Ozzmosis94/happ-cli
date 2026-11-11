// Package config manages happ-cli's YAML configuration file.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Connect holds default values for connect command flags.
type Connect struct {
	Mode          string `yaml:"mode"`
	SocksPort     int    `yaml:"socks-port"`
	HTTPPort      int    `yaml:"http-port"`
	HealthCheck   bool   `yaml:"health-check"`
	CheckInterval int    `yaml:"check-interval"`
	CheckURL      string `yaml:"check-url"`
	DNSProxy      bool   `yaml:"dns-proxy"`
	NoRouting     bool   `yaml:"no-routing"`
	SkipFirewall  bool   `yaml:"skip-firewall"`
	SystemProxy   bool   `yaml:"system-proxy"`
	Listen        string `yaml:"listen"`
	AllowLAN      bool   `yaml:"allow-lan"`
	AllowCIDRs    CIDRs  `yaml:"allow-cidr,omitempty"`
}

// CIDRs accepts both the legacy single string and a YAML sequence.
// New configurations always use a list, while old config.yaml files still load.
type CIDRs []string

func (c *CIDRs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	switch v := raw.(type) {
	case nil:
		*c = nil
	case string:
		if v == "" {
			*c = nil
		} else {
			*c = CIDRs{v}
		}
	case []interface{}:
		list := make(CIDRs, 0, len(v))
		for _, entry := range v {
			s, ok := entry.(string)
			if !ok {
				return fmt.Errorf("config: allow-cidr list must contain only strings")
			}
			list = append(list, s)
		}
		*c = list
	default:
		return fmt.Errorf("config: allow-cidr must be a string or a list of strings")
	}
	return nil
}

// File is the root config structure.
type File struct {
	Connect Connect `yaml:"connect"`
}

// Default returns a File with sensible defaults.
func Default() *File {
	return &File{
		Connect: Connect{
			Mode:          "proxy",
			Listen:        "127.0.0.1",
			SocksPort:     10808,
			HTTPPort:      10809,
			CheckInterval: 60,
			CheckURL:      "https://cloudflare.com/cdn-cgi/trace",
			DNSProxy:      true,
		},
	}
}

// Load reads config from dir/config.yaml. Returns default if the file doesn't exist.
func Load(dir string) (*File, error) {
	path := filepath.Join(dir, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}
	var c File
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", path, err)
	}
	return &c, nil
}

// Save writes the config to dir/config.yaml.
func (f *File) Save(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("config: mkdir: %w", err)
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		return fmt.Errorf("config: marshal: %w", err)
	}
	path := filepath.Join(dir, "config.yaml")
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("config: write: %w", err)
	}
	return os.Rename(tmp, path)
}
