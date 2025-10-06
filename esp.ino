// Sertakan library yang dibutuhkan
#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>

// Definisikan pin LED secara manual jika tidak ada di board definition
#ifndef LED_BUILTIN
#define LED_BUILTIN 2
#endif

// Definisikan pin untuk setiap perangkat keras
const int FAN_PIN = 25;
const int LAMP_PIN = 26;
const int LAMP_2_PIN = 27;
const int WATER_PIN = 33;

// --- KONFIGURASI ---
const char* ssid = "ARANUS_TECH";
const char* password = "Aranus140324";
const char* mqtt_server = "mqttserver.rafflesiaagro.com";
const int mqtt_port = 1883;
const char* mqtt_user = "q9HK9qUUE5u4vvRV";
const char* mqtt_pass = "TLvDMNxGYzp29DnpOcUuNR6BPhl4yGLo";
const char* coop_id = "123";

// --- TOPIK MQTT ---
char telemetry_topic[100];
char state_topic[100];

WiFiClient espClient;
PubSubClient client(espClient);
long lastMsgTime = 0;

// --- FUNGSI-FUNGSI ---

void setup_wifi() {
  delay(10);
  Serial.println();
  Serial.print("Menghubungkan ke ");
  Serial.println(ssid);
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.println("WiFi terhubung");
  Serial.print("Alamat IP: ");
  Serial.println(WiFi.localIP());
}

void controlHardware(const JsonDocument& doc) {
  if (!doc.containsKey("hardware") || !doc.containsKey("state")) {
    Serial.println("Error: Pesan tidak mengandung field 'hardware' atau 'state'.");
    return;
  }

  const char* hardware = doc["hardware"];
  bool state = doc["state"];
  int pinToControl = -1;

  if (strcmp(hardware, "fan") == 0) pinToControl = FAN_PIN;
  else if (strcmp(hardware, "lamp") == 0) pinToControl = LAMP_PIN;
  else if (strcmp(hardware, "lamp_2") == 0) pinToControl = LAMP_2_PIN;
  else if (strcmp(hardware, "water") == 0) pinToControl = WATER_PIN;

  if (pinToControl != -1) {
    digitalWrite(pinToControl, state ? HIGH : LOW);
  } else {
    Serial.print("Peringatan: Perangkat keras '");
    Serial.print(hardware);
    Serial.println("' tidak dikenal. Tidak ada aksi fisik yang dilakukan.");
    return;
  }

  Serial.print("Aksi: ");
  Serial.print(hardware);
  Serial.print(" diatur ke ");
  Serial.println(state ? "NYALA" : "MATI");

  if (!doc["meta"].isNull()) {
    // DIUBAH: Gunakan JsonObjectConst karena 'doc' adalah const
    JsonObjectConst meta = doc["meta"].as<JsonObjectConst>();
    
    if (meta.containsKey("mode")) {
      Serial.print(" > Mode: ");
      Serial.println(meta["mode"].as<const char*>());
    }
    
    if (meta.containsKey("temperature_min") || meta.containsKey("temperature_max")) {
      if (meta.containsKey("temperature_min")) {
        Serial.print(" > Batas Suhu Min: ");
        Serial.println(meta["temperature_min"].as<float>());
      }
      if (meta.containsKey("temperature_max")) {
        Serial.print(" > Batas Suhu Max: ");
        Serial.println(meta["temperature_max"].as<float>());
      }
    }

    if (meta.containsKey("schedules")) {
      // DIUBAH: Gunakan JsonArrayConst karena 'meta' adalah const
      JsonArrayConst schedules = meta["schedules"].as<JsonArrayConst>();
      Serial.println(" > Jadwal:");
      // DIUBAH: Gunakan JsonObjectConst di dalam loop
      for (JsonObjectConst schedule : schedules) {
        Serial.print("   - Order ");
        Serial.print(schedule["order"].as<int>());
        Serial.print(": ON jam ");
        Serial.print(schedule["time_on"].as<const char*>());
        Serial.print(", OFF jam ");
        Serial.println(schedule["time_off"].as<const char*>());
      }
    }
  } else {
    Serial.println(" > Tanpa metadata tambahan.");
  }
}

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.println("---");
  Serial.print("Pesan diterima pada topik: ");
  Serial.println(topic);

  JsonDocument doc;
  DeserializationError error = deserializeJson(doc, payload, length);
  if (error) {
    Serial.print("deserializeJson() gagal: ");
    Serial.println(error.c_str());
    return;
  }

  controlHardware(doc);

  Serial.println("--- Selesai Memproses Aksi ---");
}


void reconnect() {
  while (!client.connected()) {
    Serial.print("Mencoba koneksi MQTT...");
    if (client.connect(coop_id, mqtt_user, mqtt_pass)) {
      Serial.println("terhubung!");
      client.subscribe(state_topic);
      Serial.print("Berlangganan (subscribe) ke topik state: ");
      Serial.println(state_topic);
    } else {
      Serial.print("gagal, rc=");
      Serial.print(client.state());
      Serial.println(" coba lagi dalam 5 detik");
      delay(5000);
    }
  }
}

void publishTelemetry() {
  float temperature = random(2800, 3500) / 100.0;
  float humidity = random(6000, 8500) / 100.0;
  float light = random(500, 1500);

  JsonDocument doc;
  doc["temperature"] = temperature;
  doc["humidity"] = humidity;
  doc["light"] = light;

  char buffer[256];
  serializeJson(doc, buffer);
  client.publish(telemetry_topic, buffer);
  
  Serial.print("Mempublikasikan telemetri: ");
  Serial.println(buffer);
}

void setup() {
  Serial.begin(115200);
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, HIGH);
  pinMode(FAN_PIN, OUTPUT);
  pinMode(LAMP_PIN, OUTPUT);
  pinMode(LAMP_2_PIN, OUTPUT);
  pinMode(WATER_PIN, OUTPUT);
  digitalWrite(FAN_PIN, LOW);
  digitalWrite(LAMP_PIN, LOW);
  digitalWrite(LAMP_2_PIN, LOW);
  digitalWrite(WATER_PIN, LOW);
  setup_wifi();
  snprintf(telemetry_topic, sizeof(telemetry_topic), "coops/%s/telemetry", coop_id);
  snprintf(state_topic, sizeof(telemetry_topic), "coops/%s/state", coop_id);
  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);
  client.setBufferSize(512);
  digitalWrite(LED_BUILTIN, LOW);
}

void loop() {
  if (!client.connected()) {
    reconnect();
  }
  client.loop();
  long now = millis();
  if (now - lastMsgTime > 10000) {
    lastMsgTime = now;
    publishTelemetry();
  }
}