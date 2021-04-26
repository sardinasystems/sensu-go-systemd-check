package main

import (
	"fmt"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"go.uber.org/multierr"

	"github.com/sardinasystems/sensu-go-systemd-check/service"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	UnitPatterns         []string
	ExpectedActiveStates []string
	ExpectedSubStates    []string
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
		{
			Path:      "unit",
			Env:       "SYSTEMD_UNIT",
			Argument:  "unit",
			Shorthand: "s",
			Usage:     "Systemd unit(s) pattern to check",
			Value:     &plugin.UnitPatterns,
		},
		{
			Path:      "active_state",
			Env:       "SYSTEMD_ACTIVE_STATE",
			Argument:  "active",
			Shorthand: "a",
			Usage:     "Expected systemd unit(s) active state(s)",
			Value:     &plugin.ExpectedActiveStates,
			Default:   []string{"active"},
		},
		{
			Path:      "sub_state",
			Env:       "SYSTEMD_SUB_STATE",
			Argument:  "sub",
			Shorthand: "b",
			Usage:     "Expected systemd unit(s) sub state(s)",
			Value:     &plugin.ExpectedSubStates,
			Default:   []string{"running"},
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

func stringsContains(sl []string, s string) bool {
	for _, ss := range sl {
		if s == ss {
			return true
		}
	}

	return false
}

func executeCheck(event *types.Event) (int, error) {
	conn, err := dbus.New()
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("could not connect to systemd dbus: %w", err)
	}
	defer conn.Close()

	unitFetcher, err := service.InstrospectForUnitMethods(nil)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("could not introspect systemd dbus: %w", err)
	}

	unitStats, err := unitFetcher(conn, nil, plugin.UnitPatterns)
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("list units error: %w", err)
	}

	if len(unitStats) < len(plugin.UnitPatterns) {
		for _, unit := range plugin.UnitPatterns {
			matched, err := service.MatchUnitPatterns([]string{unit}, unitStats)
			if err != nil {
				fmt.Printf("CRITICAL: %s: match error: %v\n", unit, err)
				//err = multierr.Append(err, fmt.Errorf("%s: match error: %w", unit, err))
			}
			if len(matched) == 0 {
				fmt.Printf("CRITICAL: %s: not present\n", unit)
				//err = multierr.Append(err, fmt.Errorf("%s: not present", unit))
			}
		}

		//return sensu.CheckStateCritical, err
		return sensu.CheckStateCritical, nil
	}

	err = nil
	for _, unit := range unitStats {
		if !stringsContains(plugin.ExpectedActiveStates, unit.ActiveState) {
			fmt.Printf("CRITICAL: %s: active: %s\n", unit.Name, unit.ActiveState)
			err = multierr.Append(err, fmt.Errorf("%s: active: %s", unit.Name, unit.ActiveState))
			continue
		}
		if !stringsContains(plugin.ExpectedSubStates, unit.SubState) {
			fmt.Printf("CRITICAL: %s: sub: %s\n", unit.Name, unit.SubState)
			err = multierr.Append(err, fmt.Errorf("%s: sub: %s", unit.Name, unit.SubState))
			continue
		}

		fmt.Printf("OK: %s: %s and %s\n", unit.Name, unit.ActiveState, unit.SubState)
	}
	if err != nil {
		//return sensu.CheckStateCritical, err
		return sensu.CheckStateCritical, nil
	}

	return sensu.CheckStateOK, nil
}
