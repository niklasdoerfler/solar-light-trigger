package configuration

type MqttConfiguration struct {
	BrokerAddress         string
	BrokerPort            int
	Username              string
	Password              string
	ClientId              string
	TopicSolarRadiation   string
	TopicPrefixLightState string
}
