#include "SensorData.h"
#include "Global.h"
#include "WaypointManager.h"
#include <math.h>
#include <vector>
#include <HTTPClient.h>
#include "DataAnalysis.h"
#include <WiFiUdp.h>
#include <NTPClient.h>
#include <ArduinoJson.h>
#include <DHTEsp.h>

#define DHTPIN 7
#define DHTTYPE DHT22
DHT dht(DHTPIN, DHTTYPE);

WiFiUDP ntpUDP;
NTPClient timeClient(ntpUDP, "pool.ntp.org", 0, 60000);

float temperature;
float humidity;
float windSpeed;
float pressure;

// Get sensor data from DHT sensor
void getSensorData() {
  temperature = dht.readTemperature();
  humidity = dht.readHumidity();
  windSpeed = dht.readWindSpeed();
  pressure = dht.readPressure();
}

// Business logic for comparing device data with other waypoints
// deviceData: vector of SensorData objects
// returns: 0 if no anomalies detected, 1 if device issue detected, 2 if bad weather detected
int compareWithOtherWaypoints(const std::vector<SensorData>& deviceData) {
  HTTPClient http;
  String url = String(serverAddress) + "/routes/" + String(routeID) + "/get-sensor-data";

  http.begin(url);
  int httpCode = http.GET();
  String response = http.getString();

  const size_t capacity = JSON_OBJECT_SIZE(10) + JSON_ARRAY_SIZE(20);
  DynamicJsonDocument otherData(capacity);

  deserializeJson(otherData, response);
  http.end();

  std::vector<float> otherTemperatures, otherHumidities, otherWindSpeeds, otherPressures;

  for (JsonObject dataPoint : otherData.as<JsonArray>()) {
    int otherWaypointID = dataPoint["waypoint_id"];
    if (otherWaypointID != waypointID) {
      otherTemperatures.push_back(dataPoint["temperature"]);
      otherHumidities.push_back(dataPoint["humidity"]);
      otherWindSpeeds.push_back(dataPoint["wind_speed"]);
      otherPressures.push_back(dataPoint["mean_pressure"]);
    }
  }

  if (otherTemperatures.size() == 0) {
    Serial.println("No other sd detected");
    return 0;
  }

  std::vector<float> deviceTemperatures, deviceHumidities, deviceWindSpeeds, devicePressures;
  for (const SensorData& sd : deviceData) {
    deviceTemperatures.push_back(sd.temperature);
    deviceHumidities.push_back(sd.humidity);
    deviceWindSpeeds.push_back(sd.windSpeed);
    devicePressures.push_back(sd.pressure);
  }

  float deviceCVTemp = calculateCoefficientOfVariation(deviceTemperatures);
  float deviceCVHum = calculateCoefficientOfVariation(deviceHumidities);
  float deviceCVWind = calculateCoefficientOfVariation(deviceWindSpeeds);
  float deviceCVPress = calculateCoefficientOfVariation(devicePressures);

  float otherCVTemp = calculateCoefficientOfVariation(otherTemperatures);
  float otherCVHum = calculateCoefficientOfVariation(otherHumidities);
  float otherCVWind = calculateCoefficientOfVariation(otherWindSpeeds);
  float otherCVPress = calculateCoefficientOfVariation(otherPressures);

  bool tempAnomaly = deviceCVTemp > otherCVTemp * 1.5;
  bool humAnomaly = deviceCVHum > otherCVHum * 1.5;
  bool windAnomaly = deviceCVWind > otherCVWind * 1.5;
  bool pressAnomaly = deviceCVPress > otherCVPress * 1.5;

  if (tempAnomaly || humAnomaly || windAnomaly || pressAnomaly) {
    if (otherCVTemp < 0.1 && otherCVHum < 0.1 && otherCVWind < 0.1 && otherCVPress < 0.1) {
      return 1;
    }
    return 2;
  }
  return 0;
}

// Business logic for analyzing sensor data
void analyzeSensorData() {
  if (!getWeatherAlerts) return;

  const int dataSize = (sendDataFrequency * 60) / 10;
  static std::vector<SensorData> sensorDataBuffer;

  SensorData newData = {
    temperature,
    humidity,
    windSpeed,
    pressure
  };
  sensorDataBuffer.push_back(newData);

  if (sensorDataBuffer.size() > dataSize) {
    sensorDataBuffer.erase(sensorDataBuffer.begin());
  }

  if (sensorDataBuffer.size() == dataSize) {
    int anomalyType = compareWithOtherWaypoints(sensorDataBuffer);

    if (anomalyType == 1) {
      sendAnomalyStatus();
    } else if (anomalyType == 2) {
      sendBadWeatherStatus();
    }
  }
}

// Sends anomaly status to server
void sendAnomalyStatus() {
  HTTPClient http;
  String url = String(serverAddress) + "/waypoints/" + String(waypointID);

  Serial.println("Sending anomaly status (device issue) to server...");

  const size_t capacity = JSON_OBJECT_SIZE(2) + 60;
  DynamicJsonDocument doc(capacity);

  doc["status"] = "anomaly_detected";
  doc["details"] = "Device CV exceeds normal values. Possible sensor malfunction.";

  String payload;
  serializeJson(doc, payload);

  http.begin(url);
  http.addHeader("Content-Type", "application/json");
  http.addHeader("Content-Length", String(payload.length()));

  int httpCode = http.PUT(payload);
  String response = http.getString();

  Serial.println("Response: " + response);

  http.end();
}


// Sends bad weather status to server
void sendBadWeatherStatus() {
  HTTPClient http;
  String url = String(serverAddress) + "/routes/" + String(routeID);

  Serial.println("Sending bad weather status to server...");

  const size_t capacity = JSON_OBJECT_SIZE(2) + 60;
  DynamicJsonDocument doc(capacity);

  doc["status"] = "bad_weather_detected";
  doc["details"] = "CV indicates abnormal weather conditions across multiple waypoints.";

  String payload;
  serializeJson(doc, payload);

  http.begin(url);
  http.addHeader("Content-Type", "application/json");
  http.addHeader("Content-Length", String(payload.length()));

  int httpCode = http.PUT(payload);
  String response = http.getString();
  Serial.println("Response: " + response);

  http.end();
}

time_t getCurrentTime() {
  timeClient.update();
  return timeClient.getEpochTime();
}

String getFormattedTimestamp() {
  time_t now = getCurrentTime();
  struct tm* timeInfo = localtime(&now);

 
  int timezoneOffset = 3 * 3600;
  time_t localTime = now + timezoneOffset;
  timeInfo = gmtime(&localTime);

 
  char buffer[30];
  snprintf(buffer, sizeof(buffer), "%04d-%02d-%02dT%02d:%02d:%02dZ",
           timeInfo->tm_year + 1900, timeInfo->tm_mon + 1, timeInfo->tm_mday,
           timeInfo->tm_hour, timeInfo->tm_min, timeInfo->tm_sec);

  return String(buffer);
}

// Sends sensor data to server
void sendSensorData(SensorData sensorData) {
  HTTPClient http;
  String url = String(serverAddress) + "/sensor-data/";

  String timestamp = getFormattedTimestamp();

  const size_t capacity = JSON_OBJECT_SIZE(6) + 60;
  DynamicJsonDocument doc(capacity);

  doc["date"] = timestamp;
  doc["temperature"] = sensorData.temperature;
  doc["humidity"] = sensorData.humidity;
  doc["wind_speed"] = sensorData.windSpeed;
  doc["pressure"] = sensorData.pressure;
  doc["waypoint_id"] = waypointID;

  String payload;
  serializeJson(doc, payload);

  http.begin(url);
  http.addHeader("Content-Type", "application/json");
  http.addHeader("Content-Length", String(payload.length()));

  int httpCode = http.POST(payload);
  String response = http.getString();

  Serial.println("Sensor data sent: " + payload);
  Serial.println("Response: " + response);

  http.end();
}