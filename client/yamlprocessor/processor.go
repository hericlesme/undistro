/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package yamlprocessor

// Processor defines the methods necessary for creating a specific yaml
// processor.
type Processor interface {
	// GetTemplateName returns the name of the template that needs to be
	// retrieved from the source.
	GetTemplateName(version, flavor string) string

	// GetVariables parses the template blob of bytes and provides a
	// list of variables that the template requires.
	GetVariables([]byte) ([]string, error)

	// Process processes the template blob of bytes and will return the final
	// yaml with values retrieved from the values getter
	Process([]byte, func(string) (string, error)) ([]byte, error)
}
