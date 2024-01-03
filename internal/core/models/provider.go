package models

import (
	"encoding/json"
)

type ProviderType string

const (
	JsonLogFileProvider ProviderType = "jsonLogFile"
	WebSocketProvider   ProviderType = "webSocketServer"
	AggregateProvider   ProviderType = "aggregate"
)

type JsonLogFileProviderConfig struct {
	LogFile string `json:"logfile" yaml:"logFile"`
}

type WebSocketProviderConfig struct {
	BindAddress string `json:"bindAddress" yaml:"bindAddress"`
}

type AggregateProviderConfig struct {
	SubProviders []*AbstractProviderConfig `json:"subProviders" yaml:"subProviders"`
}

type AbstractProviderConfig struct {
	Name   ProviderType    `json:"name" yaml:"name"`
	Config json.RawMessage `json:"config" yaml:"config"`
}

// UnmarshalYAML - unmarshal the provider config
func (apc *AbstractProviderConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := struct {
		Provider       ProviderType `yaml:"name"`
		ProviderConfig interface{}  ` yaml:"config"`
	}{}
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}
	rawProviderConfig, err := json.Marshal(&tmp.ProviderConfig)
	if err != nil {
		return err
	}
	tmp.ProviderConfig = rawProviderConfig
	apc.Name = tmp.Provider
	apc.Config = rawProviderConfig
	return nil
}

// GetProvider - get the provider type
func (apc *AbstractProviderConfig) GetProvider() ProviderType {
	return apc.Name
}

// GetProviderConfig - unmarshal the provider config into the given interface
func (apc *AbstractProviderConfig) GetProviderConfig(config interface{}) error {
	return json.Unmarshal(apc.Config, config)
}
