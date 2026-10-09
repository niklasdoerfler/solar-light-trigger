package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// envBindings maps config keys to the environment variables that can override them.
var envBindings = map[string]string{
	"logLevel":                   "LOG_LEVEL",
	"mqtt.brokerAddress":         "MQTT_BROKER_ADDRESS",
	"mqtt.brokerPort":            "MQTT_BROKER_PORT",
	"mqtt.username":              "MQTT_USERNAME",
	"mqtt.password":              "MQTT_PASSWORD",
	"mqtt.clientId":              "MQTT_CLIENT_ID",
	"mqtt.topicSolarRadiation":   "MQTT_TOPIC_SOLAR_RADIATION",
	"mqtt.topicPrefixLightState": "MQTT_TOPIC_PREFIX_LIGHT_STATE",
}

// BindEnv registers the environment variables for all scalar config values.
func BindEnv(v *viper.Viper) error {
	for key, env := range envBindings {
		if err := v.BindEnv(key, env); err != nil {
			return err
		}
	}
	return nil
}

// ApplyTriggerEnv overrides or extends the given triggers with values from
// environment variables of the form TRIGGER_<n>_NAME, TRIGGER_<n>_ENABLED,
// TRIGGER_<n>_THRESHOLD and TRIGGER_<n>_HYSTERESIS (n starting at 0).
// Existing triggers are overridden field by field, additional indices are appended.
func ApplyTriggerEnv(triggers []TriggerConfiguration) ([]TriggerConfiguration, error) {
	for i := 0; ; i++ {
		prefix := fmt.Sprintf("TRIGGER_%d_", i)
		name, hasName := os.LookupEnv(prefix + "NAME")
		enabled, hasEnabled := os.LookupEnv(prefix + "ENABLED")
		threshold, hasThreshold := os.LookupEnv(prefix + "THRESHOLD")
		hysteresis, hasHysteresis := os.LookupEnv(prefix + "HYSTERESIS")

		if !hasName && !hasEnabled && !hasThreshold && !hasHysteresis {
			if i >= len(triggers) {
				return triggers, nil
			}
			continue
		}

		if i >= len(triggers) {
			triggers = append(triggers, TriggerConfiguration{Enabled: true})
		}
		trigger := &triggers[i]

		if hasName {
			trigger.Name = name
		}
		if hasEnabled {
			value, err := strconv.ParseBool(strings.TrimSpace(enabled))
			if err != nil {
				return nil, fmt.Errorf("invalid value for %sENABLED: %w", prefix, err)
			}
			trigger.Enabled = value
		}
		if hasThreshold {
			value, err := strconv.ParseFloat(strings.TrimSpace(threshold), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value for %sTHRESHOLD: %w", prefix, err)
			}
			trigger.Threshold = value
		}
		if hasHysteresis {
			value, err := strconv.ParseFloat(strings.TrimSpace(hysteresis), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value for %sHYSTERESIS: %w", prefix, err)
			}
			trigger.Hysteresis = value
		}
		if trigger.Name == "" {
			return nil, fmt.Errorf("trigger %d has no name, set %sNAME", i, prefix)
		}
	}
}
