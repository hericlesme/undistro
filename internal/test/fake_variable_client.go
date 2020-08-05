/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package test

import (
	"github.com/pkg/errors"
)

// FakeVariableClient provides a VariableClient backed by a map
type FakeVariableClient struct {
	variables map[string]string
}

func (f FakeVariableClient) Get(key string) (string, error) {
	if val, ok := f.variables[key]; ok {
		return val, nil
	}
	return "", errors.Errorf("value for variable %q is not set", key)
}

func (f FakeVariableClient) Set(key, value string) {
	f.variables[key] = value
}

func (f *FakeVariableClient) WithVar(key, value string) *FakeVariableClient {
	f.variables[key] = value
	return f
}

func NewFakeVariableClient() *FakeVariableClient {
	return &FakeVariableClient{
		variables: map[string]string{},
	}
}
