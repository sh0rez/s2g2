package main

import (
	_ "unsafe"

	"io"
	"net/url"

	grafana "github.com/grafana/grafana-api-golang-client"
)

//go:linkname clientRequest github.com/grafana/grafana-api-golang-client.(*Client).request
func clientRequest(c *grafana.Client, method, requestPath string, query url.Values, body io.Reader, responseStruct interface{}) error
