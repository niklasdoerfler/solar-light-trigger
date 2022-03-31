package configuration

type Configuration struct {
	LogLevel string
	Mqtt     MqttConfiguration
	Trigger  []TriggerConfiguration
}
