#include <ArduinoJson.h>
#include <vector>
#include <math.h>
#include <HTTPClient.h>
#include <WiFiUdp.h>
#include "WiFiConnect.h"
#include "WaypointManager.h"
#include "Global.h"
#include "SensorData.h"

// Device configuration before start of main loop
void setup()
{
  Serial.begin(115200);
  connectToWaypoint();
  getDeviceCoords();
  sendCoordsToServer();
  getConfigs();

  timeClient.begin();
  timeClient.update();
}

// Main loop. This loop will run indefinitely
void loop()
{
  static unsigned long lastSendTime = 0;
  unsigned long currentTime = millis();

  simulateData();

  if (currentTime - lastSendTime >= sendDataFrequency * 60 * 1000) {
    analyzeSensorData();
    lastSendTime = currentTime;

    SensorData currentData = {
      simulatedTemperature,
      simulatedHumidity,
      simulatedWindSpeed,
      simulatedPressure
    };

    sendSensorData(currentData);
  }

  delay(10 * 1000);
}
