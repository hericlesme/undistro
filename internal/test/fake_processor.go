/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package test

type FakeProcessor struct {
	errGetVariables error
	errProcess      error
	artifactName    string
}

func NewFakeProcessor() *FakeProcessor {
	return &FakeProcessor{}
}

func (fp *FakeProcessor) WithTemplateName(n string) *FakeProcessor {
	fp.artifactName = n
	return fp
}

func (fp *FakeProcessor) WithGetVariablesErr(e error) *FakeProcessor {
	fp.errGetVariables = e
	return fp
}

func (fp *FakeProcessor) WithProcessErr(e error) *FakeProcessor {
	fp.errProcess = e
	return fp
}

func (fp *FakeProcessor) GetTemplateName(version, flavor string) string {
	return fp.artifactName
}

func (fp *FakeProcessor) GetVariables(raw []byte) ([]string, error) {
	return nil, fp.errGetVariables
}

func (fp *FakeProcessor) Process(raw []byte, variablesGetter func(string) (string, error)) ([]byte, error) {
	return nil, fp.errProcess
}
