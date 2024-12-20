#include <ArduinoJson.h>
#include <HTTPClient.h>
#include "WaypointManager.h"
#include "Global.h"

unsigned int waypointID = 0;
unsigned int routeID = 0;

float x = 0.0;
float y = 0.0;

int sendDataFrequency = 20;
bool getWeatherAlerts = true;

// Connect to the server and get the waypoint and route IDs
void connectToWaypoint() {
  HTTPClient http;
  Serial.println("Getting waypoint by device serial...");

  String url = String(serverAddress) + "/waypoints/?device_serial=" + String(deviceSerial);

  http.begin(url);
  int httpCode = http.GET();
  String response = http.getString();

  Serial.println("Response : " + response);

  DynamicJsonDocument doc(1024);
  DeserializationError error = deserializeJson(doc, response);

  if (error) {
    Serial.print("JSON Deserialization failed: ");
    Serial.println(error.c_str());
    return;
  }

  if (doc.containsKey("error")) {
    Serial.print("Error: ");
    Serial.println(doc["error"].as<String>());
  } else {
    if (doc.containsKey("waypoint_id") && doc.containsKey("route_id")) {
      waypointID = doc["waypoint_id"];
      routeID = doc["route_id"];

      Serial.print("Waypoint ID: ");
      Serial.println(waypointID);

      Serial.print("Route ID: ");
      Serial.println(routeID);
    } else {
      Serial.println("Unexpected JSON response format.");
    }
  }

  http.end();
}

// Get the device coordinates
void getDeviceCoords() {
  x = random(-90, 90) + random(0, 100) / 100.0;
  y = random(-180, 180) + random(0, 100) / 100.0;
  Serial.print("Device coordinates: ");
  Serial.print("x = "); Serial.print(x);
  Serial.print(", y = "); Serial.println(y);
}

// Send the device coordinates to the server
void sendCoordsToServer() {
  HTTPClient http;
  String url = String(serverAddress) + "/waypoints/" + String(waypointID);

  Serial.println("Sending coordinates to server...");
  DynamicJsonDocument doc(256);
  doc["latitude"] = x;
  doc["longitude"] = y;

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

// Retrieve the device configurations from the server
void getConfigs() {
  HTTPClient http;
  String url = String(serverAddress) + "/device-config/" + String(waypointID);

  http.begin(url);
  int httpCode = http.GET();
  String response = http.getString();

  Serial.println("Config Response: " + response);

  DynamicJsonDocument doc(512);
  deserializeJson(doc, response);
  sendDataFrequency = doc["send_data_frequency"];
  getWeatherAlerts = doc["get_weather_alerts"];

  http.end();
}