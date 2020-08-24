/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package config

const (
	// GitHubTokenVariable defines a variable hosting the GitHub access token
	GitHubTokenVariable = "GITHUB_TOKEN"
)

// VariablesClient has methods to work with environment variables and with variables defined in the undistro configuration file.
type VariablesClient interface {
	// Get returns a variable value. If the variable is not defined an error is returned.
	// In case the same variable is defined both within the environment variables and undistro configuration file,
	// the environment variables value takes precedence.
	Get(key string) (string, error)

	// Set allows to set an explicit override for a config value.
	// e.g. It is used to set an override from a flag value over environment/config file variables.
	Set(key, values string)
}

// variablesClient implements VariablesClient.
type variablesClient struct {
	reader Reader
}

// ensure variablesClient implements VariablesClient.
var _ VariablesClient = &variablesClient{}

func newVariablesClient(reader Reader) *variablesClient {
	return &variablesClient{
		reader: reader,
	}
}

func (p *variablesClient) Get(key string) (string, error) {
	return p.reader.Get(key)
}

func (p *variablesClient) Set(key, value string) {
	p.reader.Set(key, value)
}
