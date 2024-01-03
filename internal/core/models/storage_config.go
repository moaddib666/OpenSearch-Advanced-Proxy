package models

type ConfigVersion int

const (
	ConfigVersion1 ConfigVersion = 1
)

type Config struct {
	Version ConfigVersion `yaml:"version"`
	Logs    []LogConfig   `yaml:"logs"`
}

type LogConfig struct {
	Index     string    `yaml:"index"`
	Discovery Discovery `yaml:"discovery"`
}

type DiscoveryType string

const (
	DiscoveryTypeFile DiscoveryType = "file"
	//DiscoveryTypeMock DiscoveryType = "http" // Not implemented
)

type Discovery struct {
	Type     DiscoveryType `yaml:"type"`
	Location string        `yaml:"location"`
}

type SubConfig struct {
	Version   ConfigVersion           `yaml:"version"`
	Provider  *AbstractProviderConfig `yaml:"provider"`
	Fields    map[string]*Field       `yaml:"fields"`
	Timestamp *TimestampConfig        `yaml:"timestamp"`
}

type TimestampConfig struct {
	Field string `yaml:"field"`
}
