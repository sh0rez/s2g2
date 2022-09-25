# s2g2
 
Stupid Simple Grafana GitOps

### Usage

``` bash
$ docker-compose up -d

# one of the following:
$ GRAFANA_STATE_URL='git+ssh://git@github.com/sh0rez/s2g2.git//state#main' go run .
$ GRAFANA_STATE_URL="git+file://$PWD//state#main" go run .
$ GRAFANA_STATE_URL=$PWD/state go run .
```

### State Directory

```
├── dashboards
│   └── prometheus
│       ├── overview.json
│       └── prometheus-remote-write.json
└── datasources
    └── prometheus.json
```

#### dashboards

Dashboards are organized in folders. All files matching the glob `dashboards/*/*.json` are applied.

#### datasources

All files matching the glob `datasources/*.json` are applied.

