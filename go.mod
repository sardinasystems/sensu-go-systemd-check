module github.com/sardinasystems/sensu-go-systemd-check

go 1.15

replace go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20210226220824-aa7126864d82

require (
	github.com/coreos/go-systemd/v22 v22.3.2
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/godbus/dbus/v5 v5.0.4
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.9.1 // indirect
	github.com/robertkrimen/otto v0.0.0-20200922221731-ef014fd054ac // indirect
	github.com/sensu/sensu-go/api/core/v2 v2.8.0 // indirect
	github.com/sensu/sensu-go/types v0.6.0
	github.com/sensu/sensu-plugin-sdk v0.13.1
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.3 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	go.uber.org/multierr v1.7.0
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	google.golang.org/genproto v0.0.0-20210524171403-669157292da3 // indirect
	google.golang.org/grpc v1.38.0 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)
