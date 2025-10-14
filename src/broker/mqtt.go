package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// GetJakartaTime returns current time in Asia/Jakarta timezone
func GetJakartaTime() time.Time {
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		slog.Error("Failed to load Asia/Jakarta timezone, using UTC", "error", err)
		return time.Now().UTC()
	}
	return time.Now().In(jakartaLocation)
}

// TelemetryMessage is a transient struct before being saved to the job queue.
type TelemetryMessage struct {
	CoopID    string
	Data      map[string]interface{}
	Timestamp time.Time
}

// DataStore is an interface that the database store must satisfy.
type DataStore interface {
	Save(msg TelemetryMessage, sensorTypes map[string]int32) error
}

// DIUBAH: Definisikan dan ekspor struct Config di sini.
type MQTTConfig struct {
	MQTTURL      string
	MQTTUsername string
	MQTTPassword string
}

// DIUBAH: Ubah parameter dari 'interface{}' menjadi 'MQTTConfig' yang spesifik.
func SetupMQTTClient(store DataStore, sensorTypes map[string]int32, cfg MQTTConfig) mqtt.Client {
	// DIHAPUS: Definisi 'type MQTTConfig' dan konversi 'cfg.(MQTTConfig)' tidak lagi diperlukan.

	mqttURL := cfg.MQTTURL
	mqttUser := cfg.MQTTUsername
	mqttPass := cfg.MQTTPassword

	telemetryTopic := "$share/backend-workers/coops/+/telemetry"

	opts := mqtt.NewClientOptions().AddBroker(mqttURL).SetUsername(mqttUser).SetPassword(mqttPass)
	opts.SetClientID(fmt.Sprintf("go-server-%d", time.Now().UnixNano()))
	opts.SetAutoReconnect(true).SetConnectRetry(true).SetCleanSession(true)

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		slog.Info("✅ MQTT client connected", "url", mqttURL)
		if token := c.Subscribe(telemetryTopic, 1, func(client mqtt.Client, msg mqtt.Message) {

			// DIUBAH: Logika parsing disesuaikan dengan format topik asli
			parts := strings.Split(msg.Topic(), "/")

			// DIUBAH: Pengecekan sekarang mengharapkan 3 bagian: "coops", "{id}", "telemetry"
			if len(parts) != 3 {
				slog.Warn("Received message on unexpected topic format", "topic", msg.Topic(), "expected_parts", 3, "actual_parts", len(parts))
				return
			}
			// DIUBAH: coopID sekarang berada di indeks 1
			coopID := parts[1] // coops/[0] coopID/[1] telemetry/[2]

			var sensorData map[string]interface{}
			if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
				slog.Error("Failed to unmarshal telemetry", "err", err, "topic", msg.Topic())
				return
			}

			telemetry := TelemetryMessage{
				CoopID:    coopID,
				Data:      sensorData,
				Timestamp: GetJakartaTime(),
			}

			if err := store.Save(telemetry, sensorTypes); err != nil {
				slog.Error("Failed to queue sensor data", "err", err, "coop_id", coopID)
			} else {
				slog.Info(
					"Telemetry data successfully queued for storage",
					"coop_id", coopID,
					"topic", msg.Topic(),
				)
			}
		}); token.Wait() && token.Error() != nil {
			slog.Error("Failed to subscribe to telemetry topic", "topic", telemetryTopic, "err", token.Error())
		} else {
			slog.Info("Subscribed to telemetry topic", "topic", telemetryTopic)
		}
	})

	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) { slog.Warn("❗️ MQTT connection lost", "error", err) })
	opts.SetReconnectingHandler(func(c mqtt.Client, opts *mqtt.ClientOptions) {
		slog.Info("⏳ Attempting to reconnect to MQTT broker...")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		slog.Error("❌ Initial MQTT connection failed (will retry)", "error", token.Error())
	}
	return client
}

// MqttPublish publishes a message to the MQTT broker.
func MqttPublish(c mqtt.Client, topic string, payload []byte, qos byte, retain bool) error {
	token := c.Publish(topic, qos, retain, payload)
	if token.WaitTimeout(5 * time.Second) {
		return token.Error()
	}
	return errors.New("publish timeout")
}
