package config

import (
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

type AppConfig struct {
	location string
	config   *models.Config
	indexes  map[string]*models.SubConfig
}

func (c *AppConfig) AvailableIndexes() map[string]*models.SubConfig {
	return c.indexes
}

func NewConfig(location string) *AppConfig {
	return &AppConfig{
		location: location,
		config:   &models.Config{},
		indexes:  make(map[string]*models.SubConfig),
	}
}

func (c *AppConfig) Location() string {
	return c.location
}

func (c *AppConfig) load() error {
	// load yuml config version from location
	cfgPath := path.Join(c.location, "config.yaml")
	raw, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(raw, c.config)
	if err != nil {
		return err
	}
	return nil
}

func (c *AppConfig) resolve() error {
	// resolve config version
	// resolve provider
	// resolve fields
	for _, logConfig := range c.config.Logs {
		if logConfig.Discovery.Type == models.DiscoveryTypeFile {
			subConfigPath := path.Join(logConfig.Discovery.Location, logConfig.Index+".yaml")
			subConfigData, err := os.ReadFile(subConfigPath)
			if err != nil {
				log.Errorf("Error reading sub-config file: %s", subConfigPath)
				continue
			}

			var subConfig models.SubConfig
			err = yaml.Unmarshal(subConfigData, &subConfig)
			if err != nil {
				log.Errorf("Error parsing sub-config file: %s - %s", subConfigPath, err)
				continue
			}

			log.Infof("Sub-config loaded: %s", subConfigPath)
			c.indexes[logConfig.Index] = &subConfig
		} else {
			log.Errorf("Unsupported discovery type: %s", logConfig.Discovery.Type)
			continue
		}
	}
	return nil
}

// Load loads the config from the location
func (c *AppConfig) Load() error {
	err := c.load()
	if err != nil {
		return err
	}
	err = c.resolve()
	if err != nil {
		return err
	}
	return nil
}
