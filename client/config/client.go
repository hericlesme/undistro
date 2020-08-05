/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package config

import (
	"github.com/pkg/errors"
)

// Client is used to interact with the undistro configurations.
// Clusterctl v2 handles the following configurations:
// 1. The configuration of the providers (name, type and URL of the provider repository)
// 2. Variables used when installing providers/creating clusters. Variables can be read from the environment or from the config file
// 3. The configuration about image overrides
type Client interface {
	// Providers provide access to provider configurations.
	Providers() ProvidersClient

	// Variables provide access to environment variables and/or variables defined in the undistro configuration file.
	Variables() VariablesClient

	// ImageMeta provide access to to image meta configurations.
	ImageMeta() ImageMetaClient
}

// configClient implements Client.
type configClient struct {
	reader Reader
}

// ensure configClient implements Client.
var _ Client = &configClient{}

func (c *configClient) Providers() ProvidersClient {
	return newProvidersClient(c.reader)
}

func (c *configClient) Variables() VariablesClient {
	return newVariablesClient(c.reader)
}

func (c *configClient) ImageMeta() ImageMetaClient {
	return newImageMetaClient(c.reader)
}

// Option is a configuration option supplied to New
type Option func(*configClient)

// InjectReader allows to override the default configuration reader used by undistro.
func InjectReader(reader Reader) Option {
	return func(c *configClient) {
		c.reader = reader
	}
}

// New returns a Client for interacting with the undistro configuration.
func New(path string, options ...Option) (Client, error) {
	return newConfigClient(path, options...)
}

func newConfigClient(path string, options ...Option) (*configClient, error) {
	client := &configClient{}
	for _, o := range options {
		o(client)
	}

	// if there is an injected reader, use it, otherwise use a default one
	if client.reader == nil {
		client.reader = newViperReader()
		if err := client.reader.Init(path); err != nil {
			return nil, errors.Wrap(err, "failed to initialize the configuration reader")
		}
	}

	return client, nil
}

// Reader define the behaviours of a configuration reader.
type Reader interface {
	// Init allows to initialize the configuration reader.
	Init(path string) error

	// Get returns a configuration value of type string.
	// In case the configuration value does not exists, it returns an error.
	Get(key string) (string, error)

	// Set allows to set an explicit override for a config value.
	// e.g. It is used to set an override from a flag value over environment/config file variables.
	Set(key, value string)

	// UnmarshalKey reads a configuration value and unmarshals it into the provided value object.
	UnmarshalKey(key string, value interface{}) error
}
