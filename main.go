package main

import (
	"fmt"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
	"go.uber.org/multierr"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	UnitPatterns []string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-systemd-check",
			Short:    "check for systemd service aliveness via dbus",
			Keyspace: "sensu.io/plugins/systemd/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "unit",
			Env:       "SYSTEMD_UNIT",
			Argument:  "unit",
			Shorthand: "s",
			Usage:     "Systemd unit(s) pattern to check",
			Value:     &plugin.UnitPatterns,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.UnitPatterns) == 0 {
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

	unitFetcher, err := instrospectForUnitMethods()
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("could not introspect systemd dbus: %w", err)
	}

	unitStats, err := unitFetcher(conn, nil, plugin.UnitPatterns)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("list units error: %w", err)
	}

	if len(unitStats) < len(plugin.UnitPatterns) {
		err = nil
		for _, unit := range plugin.UnitPatterns {
			matched, err := matchUnitPatterns([]string{unit}, unitStats)
			if err != nil {
				fmt.Printf("CRITICAL: %s: match error: %v\n", unit, err)
				err = multierr.Append(err, fmt.Errorf("%s: match error: %w", unit, err))
			}
			if len(matched) == 0 {
				fmt.Printf("CRITICAL: %s: not present\n", unit)
				err = multierr.Append(err, fmt.Errorf("%s: not present", unit))
			}
		}

		//return sensu.CheckStateCritical, err
		return sensu.CheckStateCritical, nil
	}

	err = nil
	for _, unit := range unitStats {
		if unit.ActiveState != "active" {
			fmt.Printf("CRITICAL: %s: active: %s\n", unit.Name, unit.ActiveState)
			err = multierr.Append(err, fmt.Errorf("%s: active: %s", unit.Name, unit.ActiveState))
			continue
		}
		if unit.SubState != "running" {
			fmt.Printf("CRITICAL: %s: sub: %s\n", unit.Name, unit.SubState)
			err = multierr.Append(err, fmt.Errorf("%s: sub: %s", unit.Name, unit.SubState))
			continue
		}

		fmt.Printf("OK: %s: active and running\n", unit.Name)
	}
	if err != nil {
		//return sensu.CheckStateCritical, err
		return sensu.CheckStateCritical, nil
	}

	return sensu.CheckStateOK, nil
}
