package main

import (
	"fmt"
	"log"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
	"go.uber.org/multierr"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Units []string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-systemd-check",
			Short:    "check for systemd service aliveness via dbus",
			Keyspace: "sensu.io/plugins/sensu-go-systemd-check/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "unit",
			Env:       "SYSTEMD_UNIT",
			Argument:  "unit",
			Shorthand: "s",
			Usage:     "Systemd unit(s) to check",
			Value:     &plugin.Units,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func ssContains(s string, ss []string) bool {
	for _, sl := range ss {
		if s == sl {
			return true
		}
	}

	return false
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.Units) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--unit or SYSTEMD_UNIT environment variable is required")
	}

	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	conn, err := dbus.New()
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("could not connect to systemd dbus: %w", err)
	}
	defer conn.Close()

	unitStats, err := conn.ListUnitsByNames(plugin.Units)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("list units error: %w", err)
	}
	if len(unitStats) < len(plugin.Units) {
		err = nil
		foundUnits := make([]string, 0, len(unitStats))
		for _, unit := range unitStats {
			foundUnits = append(foundUnits, unit.Name)
		}
		for _, unit := range plugin.Units {
			if !ssContains(unit, foundUnits) {
				err = multierr.Append(err, fmt.Errorf("%s: not present", unit))
			}
		}

		return sensu.CheckStateCritical, err
	}

	for _, unit := range unitStats {
		log.Printf("%s: %v\n", unit.Name, unit)

		if unit.ActiveState != "active" {
			return sensu.CheckStateCritical, fmt.Errorf("%s: active: %s", unit.Name, unit.ActiveState)
		}
		if unit.SubState != "running" {
			return sensu.CheckStateCritical, fmt.Errorf("%s: sub: %s", unit.Name, unit.SubState)
		}

		log.Printf("%s: OK\n", unit.Name)
	}

	return sensu.CheckStateOK, nil
}
