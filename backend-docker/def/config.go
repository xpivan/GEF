package def

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration keeps the configs for the entire application
type Configuration struct {
	Docker []DockerConfig
	Server ServerConfig

	// TmpDir is the directory to keep session files in
	// If the path is relative, it will be used as a subfolder of the system temporary directory
	TmpDir string
}

// DockerConfig configuration for building docker clients
type DockerConfig struct {
	UseBoot2Docker bool
	Endpoint       string
	Description    string
}

// ServerConfig keeps the configuration options needed to make a Server
type ServerConfig struct {
	Address          string
	ReadTimeoutSecs  int
	WriteTimeoutSecs int
}

func (c DockerConfig) String() string {
	if c.Endpoint != "" {
		return fmt.Sprintf("Endpoint: %s -- %s", c.Endpoint, c.Description)
	} else if c.UseBoot2Docker {
		return fmt.Sprintf("Boot2Docker[env] -- %s", c.Description)
	}
	return fmt.Sprintf("unknown -- %s", c.Description)
}

// ReadConfigFile reads a configuration file
func ReadConfigFile(configFilepath string) (Configuration, error) {
	var config Configuration

	file, err := os.Open(configFilepath)
	if err != nil {
		return config, Err(err, "Cannot open config file %s", configFilepath)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, Err(err, "Cannot read config file %s", configFilepath)
	}

	if len(config.Docker) == 0 {
		return config, Err(nil, "Docker configuration missing in file: %s", configFilepath)
	}

	return config, nil
}