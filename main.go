package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"solar-light-trigger/model/configuration"
	"strconv"
	"strings"
	"syscall"
)

var (
	BuildVersion = "dev"
	BuildTime    = "-"
)

const (
	UNDEFINED int = 0
	DAY           = 1
	NIGHT         = 2
)

var (
	mqttClient      mqtt.Client
	config          configuration.Configuration
	lightStates     []int
	lightStateNames = [3]string{"undefined", "day", "night"}
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Debugf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var messageSolarRadiationHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Debugf("Solar radiation message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	var value = parseFloatFromMessage(msg)
	checkLightState(value)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("loglevel", "info")

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Error reading config file, using default values. %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logLevel, err := log.ParseLevel(config.LogLevel)
	if err == nil {
		log.SetLevel(logLevel)
	}

	lightStates = make([]int, len(config.Trigger))
	for i, trigger := range config.Trigger {
		lightStates[i] = UNDEFINED
		log.Infof("Configuring trigger: '%s', Enabled: %t, Threshold: %f, Hysteresis: %f", trigger.Name, trigger.Enabled, trigger.Threshold, trigger.Hysteresis)
	}

	log.Debug("Config loaded.")
}

func parseFloatFromMessage(msg mqtt.Message) float64 {
	float, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		log.Error("Unable to parse float from mqtt payload:", err)
		return 0.0
	}
	return float
}

func checkLightState(currentValue float64) {
	for i, trigger := range config.Trigger {
		if trigger.Enabled {
			if (lightStates[i] == DAY || lightStates[i] == UNDEFINED) && currentValue < (trigger.Threshold-trigger.Hysteresis) {
				changeLightState(i, NIGHT, currentValue)
			} else if (lightStates[i] == NIGHT || lightStates[i] == UNDEFINED) && currentValue > (trigger.Threshold+trigger.Hysteresis) {
				changeLightState(i, DAY, currentValue)
			} else if lightStates[i] == UNDEFINED && currentValue > trigger.Threshold-trigger.Hysteresis {
				changeLightState(i, DAY, currentValue)
			}
		}
	}
}

func changeLightState(i int, lightState int, currentValue float64) {
	lightStates[i] = lightState
	trigger := config.Trigger[i]
	log.Infof("Trigger '%s' changed to '%s' (triggered by light level: %f)", trigger.Name, lightStateNames[lightStates[i]], currentValue)
	publishMessage(fmt.Sprintf("%s/%s", strings.TrimSuffix(config.Mqtt.TopicPrefixLightState, "/"), trigger.Name), lightStateNames[lightStates[i]])
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	rOps := client.OptionsReader()
	servers := rOps.Servers()
	log.Infof("Connected to mqtt broker: %s:%s", servers[0].Hostname(), servers[0].Port())
	token := client.Subscribe(config.Mqtt.TopicSolarRadiation, 0, messageSolarRadiationHandler)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Unable to subscribe: ", token.Error())
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Warnf("Connect to mqtt broker lost: %v", err)
}

func publishMessage(topic string, message string) {
	log.Debugf("Publish message on topic %s: '%s'", topic, message)
	token := mqttClient.Publish(topic, 0, true, message)
	token.Wait()
}

func SetupMqttConnection() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Mqtt.BrokerAddress, config.Mqtt.BrokerPort))
	opts.SetClientID(config.Mqtt.ClientId)
	opts.SetUsername(config.Mqtt.Username)
	opts.SetPassword(config.Mqtt.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Error("Unable to connect to mqtt broker:", token.Error())
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	log.Infof("Hello Solar Light Trigger %s! ☀️💡 (Build: %s)", BuildVersion, BuildTime)

	loadConfig()

	keepAlive := make(chan os.Signal)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)
	SetupMqttConnection()
	<-keepAlive
}
