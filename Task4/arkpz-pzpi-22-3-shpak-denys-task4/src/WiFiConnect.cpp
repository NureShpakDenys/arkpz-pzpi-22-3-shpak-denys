#include "WiFiConnect.h"

const char *ssid = "ESP32-Access-Point";
const char *password = "123456789";

AsyncWebServer server(80);

// Access Point creation and HTTP server initialization
void connectToWiFi()
{
  WiFi.softAP(ssid, password);
  IPAddress IP = WiFi.softAPIP();
  Serial.println("Access Point started");
  Serial.print("AP IP address: ");
  Serial.println(IP);

  server.on("/", HTTP_GET, [](AsyncWebServerRequest *request)
            {
    String html = "<html><body>";
    html += "<h1>Wi-Fi Configuration</h1>";
    html += "<form action='/connect' method='POST'>";
    html += "SSID: <input type='text' name='ssid'><br>";
    html += "Password: <input type='password' name='password'><br>";
    html += "<input type='submit' value='Connect'>";
    html += "</form>";
    html += "</body></html>";
    request->send(200, "text/html", html); });

  server.on("/connect", HTTP_POST, [](AsyncWebServerRequest *request)
            {
    String apssid = request->arg("ssid");
    String appassword = request->arg("password");

    WiFi.mode(WIFI_STA);
    WiFi.begin(apssid.c_str(), appassword.c_str());

    int attempt = 0;
    while (WiFi.status() != WL_CONNECTED && attempt < 20) {
      delay(500);
      Serial.print(".");
      attempt++;
    }

    if (WiFi.status() == WL_CONNECTED) {

      String message = "Connected to Wi-Fi! IP: ";
      message += WiFi.localIP().toString();
      request->send(200, "text/html", message);
      Serial.println("\nConnected to Wi-Fi!");
      Serial.print("IP Address: ");
      Serial.println(WiFi.localIP());
    } else {

      WiFi.disconnect();
      WiFi.mode(WIFI_AP);
      WiFi.softAP(ssid, password);

      request->send(200, "text/html", "Failed to connect. Returning to Access Point mode.");
      Serial.println("\nFailed to connect to Wi-Fi. Returning to Access Point mode.");
    } });

  server.begin();
  Serial.println("HTTP server started");
}