/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/google/go-github/v32/github"
)

const baseURLPath = "/api-v3"

// NewFakeGitHub sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func NewFakeGitHub() (client *github.Client, mux *http.ServeMux, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is configured to use test server.
	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	return client, mux, server.Close
}
