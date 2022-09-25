package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"

	grafana "github.com/grafana/grafana-api-golang-client"
	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/hairyhenderson/go-fsimpl/httpfs"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	flag.Parse()
	grafanaUrl := flag.Arg(0)
	if grafanaUrl == "" {
		grafanaUrl = "http://admin:admin@127.0.0.1:3000"
	}
	api, err := grafana.New(grafanaUrl, grafana.Config{})
	if err != nil {
		return err
	}

	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	mux.Add(httpfs.FS)
	mux.Add(gitfs.FS)

	stateStr := os.Getenv("GRAFANA_STATE_URL")
	if stateStr == "" {
		return fmt.Errorf("GRAFANA_STATE_URL is required")
	}

	src, err := mux.Lookup(stateStr)
	if err != nil {
		return err
	}

	// test filesystem connection
	_, err = fs.ReadDir(src, ".")
	if err != nil {
		return err
	}

	syncers := []Syncer{
		syncer[Dashboard]{NewDashboardProvider(api)},
		syncer[Datasource]{NewDatasourceProvider(api)},
	}

	for _, s := range syncers {
		go s.Sync(src, time.NewTicker(1*time.Minute).C)
	}

	select {}
}

type Trigger <-chan time.Time

type Syncer interface {
	Sync(src fs.FS, t Trigger) error
}

type syncer[T any] struct {
	Provider[T]
}

func (s syncer[T]) Sync(src fs.FS, t Trigger) error {
	for ; true; <-t {
		ts, err := s.Provider.Parse(src)
		if err != nil {
			return err
		}

		ok := 0
		for _, t := range ts {
			if err := s.Provider.Store(t); err != nil {
				log.Println(err)
				continue
			}
			ok++
		}
		log.Printf("%T: %d/%d ok", *new(T), ok, len(ts))
	}

	return nil
}

type Provider[T any] interface {
	Store(T) error
	Parse(fs.FS) ([]T, error)
}
