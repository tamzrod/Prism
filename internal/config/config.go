// internal/config/config.go

package config

import (
    "os"
    "strconv"

    "gopkg.in/yaml.v3"
)

// Config holds the minimal configuration for the PRISM server.
// More fields can be added as needed, but for now only the listening
// port is required.

type Config struct {
    // Port the server should listen on (TCP).
    Port int `json:"port"`
}

// Load reads configuration from the given YAML file path. If the file does not
// exist or cannot be parsed the returned error describes the failure.
//
// If path is empty the function will look for an environment variable
// PRISM_PORT and use that value (parsed as an integer). If neither the file
// nor the env var is provided, a default port of 12345 is returned.
//
// The configuration file is expected to be YAML with a top-level `port`
// field. Example:
//
//   port: 4321
func Load(path string) (*Config, error) {
    if path != "" {
        f, err := os.Open(path)
        if err != nil {
            return nil, err
        }
        defer f.Close()
        var cfg Config
        if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
            return nil, err
        }
        if cfg.Port == 0 {
            cfg.Port = 12345
        }
        return &cfg, nil
    }
    // no file path, check env var
    if p := os.Getenv("PRISM_PORT"); p != "" {
        var cfg Config
        // simple atoi
        var err error
        cfg.Port, err = strconv.Atoi(p)
        if err != nil {
            return nil, err
        }
        return &cfg, nil
    }
    // default
    return &Config{Port: 12345}, nil
}
