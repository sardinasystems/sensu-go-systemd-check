module github.com/sardinasystems/sensu-go-systemd-check

go 1.13

require (
	github.com/coreos/etcd v3.3.20+incompatible // indirect
	github.com/coreos/go-systemd/v22 v22.0.0
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/godbus/dbus v4.1.0+incompatible
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/mitchellh/mapstructure v1.2.2 // indirect
	github.com/pelletier/go-toml v1.7.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sensu-community/sensu-plugin-sdk v0.6.0
	github.com/sensu/sensu-go v0.0.0-20200131164840-40b1d5938251
	github.com/sirupsen/logrus v1.5.0 // indirect
	github.com/spf13/cobra v0.0.7 // indirect
	go.uber.org/multierr v1.5.0
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	google.golang.org/genproto v0.0.0-20200403120447-c50568487044 // indirect
	google.golang.org/grpc v1.28.0 // indirect
	gopkg.in/ini.v1 v1.55.0 // indirect
)

//replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v0.0.0-20200316104309-cb8b64719ae3
