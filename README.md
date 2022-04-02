# Solar Light Trigger

This small tool can be used to evaluate measured solar radiation values to determine whether it is day or night.
Therefore, multiple trigger can be specified, where each of them has its own threshold and hysteresis config.

Data exchange is realised via a MQTT connection. This service subscribes to the MQTT topic `topicSolarRadiation` which
provides some kind of brightness value (e.g. solar radiation measured in watts per square meter). Based on the trigger
level the service publishes a message for each trigger changing state from day to night and vice versa on the topic
named by the prefix `topicPrefixLightState` and the trigger name (see config below).

## Setup

Simply compile this golang project or download a precompiled binary for your platform from the release page.

## Config

The service can be configured by a config.yaml file placed next to the binary.

The following represents an example config:

```yaml
logLevel: info

mqtt:
  brokerAddress: 1.2.3.4
  brokerPort: 1883
  username: username
  password: password
  clientId: solar-light-trigger
  topicSolarRadiation: solarRadiation
  topicPrefixLightState: lightState/

trigger:
  - name: indoor
    enabled: true
    threshold: 10.0
    hysteresis: 1.0

  - name: outdoor
    enabled: false
    threshold: 2.5
    hysteresis: 1.0
```
