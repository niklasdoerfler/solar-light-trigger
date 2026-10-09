# Solar Light Trigger

This small tool can be used to evaluate measured solar radiation values to determine whether it is day or night.
Therefore, multiple trigger can be specified, where each of them has its own threshold and hysteresis config.

Data exchange is realised via a MQTT connection. This service subscribes to the MQTT topic `topicSolarRadiation` which
provides some kind of brightness value (e.g. solar radiation measured in watts per square meter). Based on the trigger
level the service publishes a message for each trigger changing state from day to night and vice versa on the topic
named by the prefix `topicPrefixLightState` and the trigger name (see config below).

## Setup

Simply compile this golang project or download a precompiled binary for your platform from the release page.

### Container

A multi-arch container image (`linux/amd64`, `linux/arm64`, `linux/arm/v7`) is published to the GitHub Container Registry:

```shell
docker run -d --name solar-light-trigger --restart unless-stopped \
  -e MQTT_BROKER_ADDRESS=1.2.3.4 \
  -e MQTT_USERNAME=username \
  -e MQTT_PASSWORD=password \
  -e TRIGGER_0_NAME=indoor \
  -e TRIGGER_0_THRESHOLD=10.0 \
  -e TRIGGER_0_HYSTERESIS=1.0 \
  ghcr.io/niklasdoerfler/solar-light-trigger:latest
```

Alternatively mount a config file to `/config/config.yaml`:

```shell
docker run -d --name solar-light-trigger --restart unless-stopped \
  -v ./config.yaml:/config/config.yaml:ro \
  ghcr.io/niklasdoerfler/solar-light-trigger:latest
```

Or with Docker Compose:

```yaml
services:
  solar-light-trigger:
    image: ghcr.io/niklasdoerfler/solar-light-trigger:latest
    restart: unless-stopped
    environment:
      MQTT_BROKER_ADDRESS: 1.2.3.4
      MQTT_USERNAME: username
      MQTT_PASSWORD: password
      TRIGGER_0_NAME: indoor
      TRIGGER_0_THRESHOLD: 10.0
      TRIGGER_0_HYSTERESIS: 1.0
      TRIGGER_1_NAME: outdoor
      TRIGGER_1_ENABLED: false
      TRIGGER_1_THRESHOLD: 2.5
      TRIGGER_1_HYSTERESIS: 1.0
```

Available tags:

| Tag                       | Description                       |
|---------------------------|-----------------------------------|
| `latest`                  | Latest release                    |
| `1.2.3`, `1.2`, `1`, `v1.2.3` | Specific release versions |
| `edge`                    | Latest build of the `main` branch |
| `sha-<commit>`            | Build of a specific commit        |

To build the image yourself:

```shell
docker build -f Containerfile -t solar-light-trigger .
# or
podman build -t solar-light-trigger .
```

## Config

The service can be configured by a `config.yaml` file placed next to the binary (or in `/config`), by environment
variables, or a combination of both. Environment variables take precedence over values from the config file.

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

### Environment variables

| Variable                        | Config key                   | Default               |
|---------------------------------|------------------------------|-----------------------|
| `LOG_LEVEL`                     | `logLevel`                   | `info`                |
| `MQTT_BROKER_ADDRESS`           | `mqtt.brokerAddress`         |                       |
| `MQTT_BROKER_PORT`              | `mqtt.brokerPort`            | `1883`                |
| `MQTT_USERNAME`                 | `mqtt.username`              |                       |
| `MQTT_PASSWORD`                 | `mqtt.password`              |                       |
| `MQTT_CLIENT_ID`                | `mqtt.clientId`              | `solar-light-trigger` |
| `MQTT_TOPIC_SOLAR_RADIATION`    | `mqtt.topicSolarRadiation`   |                       |
| `MQTT_TOPIC_PREFIX_LIGHT_STATE` | `mqtt.topicPrefixLightState` |                       |

Triggers are configured with indexed variables, starting at `0` and counting up without gaps:

| Variable                 | Config key             | Default |
|--------------------------|------------------------|---------|
| `TRIGGER_<n>_NAME`       | `trigger[n].name`      |         |
| `TRIGGER_<n>_ENABLED`    | `trigger[n].enabled`   | `true`  |
| `TRIGGER_<n>_THRESHOLD`  | `trigger[n].threshold` | `0`     |
| `TRIGGER_<n>_HYSTERESIS` | `trigger[n].hysteresis` | `0`     |

If a trigger with index `n` already exists in the config file, only the fields set via environment variables are
overridden. Otherwise, a new trigger is added, which requires `TRIGGER_<n>_NAME` to be set.
