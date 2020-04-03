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
	Units           []string
	UseByNameSearch bool
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
			Usage:     "Systemd unit(s) to check",
			Value:     &plugin.Units,
		},
		&sensu.PluginConfigOption{
			Path:      "use_by_name_search",
			Env:       "SYSTEMD_USE_BY_NAME_SEARCH",
			Argument:  "by-name",
			Shorthand: "n",
			Usage:     "use by name search (use with systemd > 230)",
			Value:     &plugin.UseByNameSearch,
			Default:   false,
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

	var unitStats []dbus.UnitStatus

	if plugin.UseByNameSearch {
		unitStats, err = conn.ListUnitsByNames(plugin.Units)
	} else {
		unitStats, err = conn.ListUnitsByPatterns(nil, plugin.Units)
	}
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
