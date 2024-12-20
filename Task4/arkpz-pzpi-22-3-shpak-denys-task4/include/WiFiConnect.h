#ifndef WIFICONNECT_H
#define WIFICONNECT_H

#include <WiFi.h>
#include <ESPAsyncWebServer.h>

extern const char* ssid;
extern const char* password;

extern AsyncWebServer server;

void connectToWiFi(const char *networkSSID, const char *networkPassword);

#endif