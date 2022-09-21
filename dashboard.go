package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	grafana "github.com/grafana/grafana-api-golang-client"
)

type Dashboard grafana.Dashboard

func NewDashboardProvider(api *grafana.Client) Provider[Dashboard] {
	return &DashboardProvider{api: api}
}

type DashboardProvider struct {
	api *grafana.Client
}

func (d *DashboardProvider) Store(db Dashboard) error {
	db.Overwrite = true

	_, err := d.api.NewDashboard(grafana.Dashboard(db))
	return err
}

func (d *DashboardProvider) Parse(src fs.FS) ([]Dashboard, error) {
	files, err := fs.Glob(src, "dashboards/*/*.json")
	if err != nil {
		panic(err)
	}

	dbs := make([]Dashboard, 0, len(files))
	for _, name := range files {
		f, err := src.Open(name)
		if err != nil {
			log.Printf("%s: %s", name, err)
			continue
		}

		var db Dashboard
		db.FolderUID = strings.TrimPrefix(filepath.Dir(name), "dashboards/")
		if err := json.NewDecoder(f).Decode(db.Model); err != nil {
			log.Printf("%s: %s", name, err)
			continue
		}
		dbs = append(dbs, db)
	}

	return dbs, nil
}
