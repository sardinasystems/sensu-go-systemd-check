module github.com/sardinasystems/sensu-go-systemd-check

go 1.13

require (
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/coreos/go-systemd/v22 v22.0.0
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/godbus/dbus v4.1.0+incompatible
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/mitchellh/mapstructure v1.3.1 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sensu-community/sensu-plugin-sdk v0.7.0
	github.com/sensu/sensu-go v0.0.0-20200131164840-40b1d5938251
	go.uber.org/multierr v1.5.0
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2 // indirect
	golang.org/x/sys v0.0.0-20200523222454-059865788121 // indirect
	google.golang.org/genproto v0.0.0-20200527145253-8367513e4ece // indirect
	google.golang.org/grpc v1.29.1 // indirect
	gopkg.in/ini.v1 v1.56.0 // indirect
)

//replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v0.0.0-20200316104309-cb8b64719ae3
