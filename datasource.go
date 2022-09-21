package main

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"log"
	"strings"

	grafana "github.com/grafana/grafana-api-golang-client"
)

type Datasource grafana.DataSource

type DatasourceProvider struct {
	api *grafana.Client
}

func NewDatasourceProvider(api *grafana.Client) Provider[Datasource] {
	return &DatasourceProvider{api: api}
}

func (d *DatasourceProvider) Store(ds Datasource) error {
	data, err := json.Marshal(ds)
	if err != nil {
		return err
	}

	// TODO: this is terrifying. make the client return a properly structured error instead
	err = clientRequest(d.api, "PUT", "/api/datasources/uid/"+ds.UID, nil, bytes.NewBuffer(data), nil)
	if err != nil && strings.Contains(err.Error(), "status: 404") {
		_, err := d.api.NewDataSource((*grafana.DataSource)(&ds))
		return err
	}

	return err
}

func (d *DatasourceProvider) Parse(src fs.FS) ([]Datasource, error) {
	files, err := fs.Glob(src, "datasources/*.json")
	if err != nil {
		return nil, err
	}

	dss := make([]Datasource, 0, len(files))
	for _, name := range files {
		f, err := src.Open(name)
		if err != nil {
			log.Printf("%s: %s", name, err)
			continue
		}

		var ds Datasource
		if err := json.NewDecoder(f).Decode(&ds); err != nil {
			log.Printf("%s: %s", name, err)
			continue
		}

		dss = append(dss, ds)
	}

	return dss, nil
}
