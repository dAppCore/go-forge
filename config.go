package forge

import (
	"encoding/json"
	"os"
	"path/filepath"

	core "dappco.re/go/core"
	coreio "dappco.re/go/core/io"
)

const (
	// DefaultURL is the fallback Forgejo instance URL when neither flag nor
	// environment variable is set.
	//
	// Usage:
	//  cfgURL, _, _ := forge.ResolveConfig("", "")
	//  _ = cfgURL == forge.DefaultURL
	DefaultURL = "http://localhost:3000"
)

const defaultConfigPath = ".config/forge/config.json"

type configFile struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", core.E("configPath", "forge: resolve home directory", err)
	}
	return filepath.Join(home, defaultConfigPath), nil
}

func readConfigFile() (url, token string, err error) {
	path, err := configPath()
	if err != nil {
		return "", "", err
	}

	data, err := coreio.Local.Read(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", nil
		}
		return "", "", core.E("ResolveConfig", "forge: read config file", err)
	}

	var cfg configFile
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		return "", "", core.E("ResolveConfig", "forge: decode config file", err)
	}
	return cfg.URL, cfg.Token, nil
}

// SaveConfig persists the Forgejo URL and API token to the default config file.
// It creates the parent directory if it does not already exist.
//
// Usage:
//
//	_ = forge.SaveConfig("https://forge.example.com", "token")
func SaveConfig(url, token string) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := coreio.Local.EnsureDir(filepath.Dir(path)); err != nil {
		return core.E("SaveConfig", "forge: create config directory", err)
	}
	payload, err := json.MarshalIndent(configFile{URL: url, Token: token}, "", "  ")
	if err != nil {
		return core.E("SaveConfig", "forge: encode config file", err)
	}
	return coreio.Local.WriteMode(path, string(payload), 0600)
}

// ResolveConfig resolves the Forgejo URL and API token from flags, environment
// variables, config file, and built-in defaults. Priority order:
// flags > env > config file > defaults.
//
// Environment variables:
//   - FORGE_URL   — base URL of the Forgejo instance
//   - FORGE_TOKEN — API token for authentication
//
// Usage:
//
//	url, token, err := forge.ResolveConfig("", "")
//	_ = url
//	_ = token
func ResolveConfig(flagURL, flagToken string) (url, token string, err error) {
	if envURL, ok := os.LookupEnv("FORGE_URL"); ok && envURL != "" {
		url = envURL
	}
	if envToken, ok := os.LookupEnv("FORGE_TOKEN"); ok && envToken != "" {
		token = envToken
	}

	if flagURL != "" {
		url = flagURL
	}
	if flagToken != "" {
		token = flagToken
	}
	if url == "" || token == "" {
		fileURL, fileToken, fileErr := readConfigFile()
		if fileErr != nil {
			return "", "", fileErr
		}
		if url == "" {
			url = fileURL
		}
		if token == "" {
			token = fileToken
		}
	}
	if url == "" {
		url = DefaultURL
	}
	return url, token, nil
}

// NewFromConfig creates a new Forge client using resolved configuration.
//
// Usage:
//
//	f, err := forge.NewFromConfig("", "")
//	_ = f
func NewFromConfig(flagURL, flagToken string, opts ...Option) (*Forge, error) {
	return NewForgeFromConfig(flagURL, flagToken, opts...)
}

// NewForgeFromConfig creates a new Forge client using resolved configuration.
// It returns an error if no API token is available from flags or environment.
//
// Usage:
//
//	f, err := forge.NewForgeFromConfig("", "")
//	_ = f
func NewForgeFromConfig(flagURL, flagToken string, opts ...Option) (*Forge, error) {
	url, token, err := ResolveConfig(flagURL, flagToken)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, core.E("NewForgeFromConfig", "forge: no API token configured (set FORGE_TOKEN or pass --token)", nil)
	}
	return NewForge(url, token, opts...), nil
}
