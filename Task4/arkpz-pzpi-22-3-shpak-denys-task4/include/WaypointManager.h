#ifndef WAYPOINT_MANAGER_H
#define WAYPOINT_MANAGER_H

#include <WiFi.h>

extern unsigned int waypointID;
extern unsigned int routeID;

extern float x;
extern float y;

extern int sendDataFrequency;
extern bool getWeatherAlerts;

// Function to connect to the waypoint
void connectToWaypoint();

// Function to get the current location of the device
void getDeviceCoords();

// Function to send the current location of the device to the server
void sendCoordsToServer();

// Function to get the configs from the server
void getConfigs();

#endif