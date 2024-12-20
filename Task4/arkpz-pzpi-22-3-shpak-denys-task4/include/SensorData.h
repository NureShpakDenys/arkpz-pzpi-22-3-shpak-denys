#ifndef SENSOR_DATA_H
#define SENSOR_DATA_H

#include <WiFiUdp.h>
#include <NTPClient.h>
#include <vector>

extern float temperature;
extern float humidity;
extern float windSpeed;
extern float pressure;

extern NTPClient timeClient;

extern struct SensorData {
  float temperature;
  float humidity;
  float windSpeed;
  float pressure;
} sensorData;

// Function to get sensor data from DHT sensor
void getSensorData();

// Function to analyze sensor data
void analyzeSensorData();

// Function to compare device data with other waypoints
int compareWithOtherWaypoints(const std::vector<SensorData>& deviceData);

// Function to send anomaly status to server
void sendAnomalyStatus();

// Function to send bad weather status to server
void sendBadWeatherStatus();

// Function to send sensor data to server
void sendSensorData(SensorData sensorData);

#endif