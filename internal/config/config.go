package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type LinkConfig struct {
	Name   string `yaml:"name"`
	Target string `yaml:"target"`
}

type VaultConfig struct {
	Description string       `yaml:"description"`
	Links       []LinkConfig `yaml:"links"`
}

type Config struct {
	Vaults map[string]VaultConfig `yaml:"vaults"`
}

var (
	configPath string
	appRoot    string
)

func SetConfigPath(path string) {
	configPath = path
	appRoot = filepath.Dir(path)
}

func GetConfigPath() string {
	return configPath
}

func GetAppRoot() string {
	return appRoot
}

func GetVaultsDir() string {
	return filepath.Join(appRoot, "vaults")
}

func ResolveConfigPath() (string, error) {
	if env := os.Getenv("VAULTEX_CONFIG"); env != "" {
		return env, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	candidate := filepath.Join(cwd, "vaultex.yaml")
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}

	exe, err := os.Executable()
	if err == nil {
		candidate = filepath.Join(filepath.Dir(exe), "vaultex.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return filepath.Join(cwd, "vaultex.yaml"), nil
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Vaults: make(map[string]VaultConfig)}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Vaults == nil {
		cfg.Vaults = make(map[string]VaultConfig)
	}

	return &cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (c *Config) AddVault(name, description string) error {
	if _, exists := c.Vaults[name]; exists {
		return fmt.Errorf("vault %q already exists", name)
	}

	c.Vaults[name] = VaultConfig{
		Description: description,
		Links:       []LinkConfig{},
	}
	return nil
}

func (c *Config) RemoveVault(name string) error {
	if _, exists := c.Vaults[name]; !exists {
		return fmt.Errorf("vault %q does not exist", name)
	}

	delete(c.Vaults, name)
	return nil
}

func (c *Config) AddLink(vaultName, linkName, target string) error {
	v, exists := c.Vaults[vaultName]
	if !exists {
		return fmt.Errorf("vault %q does not exist", vaultName)
	}

	for _, l := range v.Links {
		if l.Name == linkName {
			return fmt.Errorf("link %q already exists in vault %q", linkName, vaultName)
		}
	}

	v.Links = append(v.Links, LinkConfig{Name: linkName, Target: target})
	c.Vaults[vaultName] = v
	return nil
}

func (c *Config) RemoveLink(vaultName, linkName string) error {
	v, exists := c.Vaults[vaultName]
	if !exists {
		return fmt.Errorf("vault %q does not exist", vaultName)
	}

	for i, l := range v.Links {
		if l.Name == linkName {
			v.Links = append(v.Links[:i], v.Links[i+1:]...)
			c.Vaults[vaultName] = v
			return nil
		}
	}

	return fmt.Errorf("link %q not found in vault %q", linkName, vaultName)
}

func (c *Config) GetVault(name string) (VaultConfig, error) {
	v, exists := c.Vaults[name]
	if !exists {
		return VaultConfig{}, fmt.Errorf("vault %q does not exist", name)
	}
	return v, nil
}

func (c *Config) VaultNames() []string {
	names := make([]string, 0, len(c.Vaults))
	for name := range c.Vaults {
		names = append(names, name)
	}
	return names
}
