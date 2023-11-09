# README.md

### Welcome

**Overview**

Microxt is a multiplatform application developed upon Golang and the Electron Framework API for interfacing with IoT devices. The application scans for devices on the local area network operating on UDP out of port 8040 and connects to devices that respond properly to a handshake request which requires an public/private key match.

Microxt handles the request/response routing via key based messages, which specify content type or response sequences to follow. 



### Device Network Notes

We expect the following network connection parameters:

1. Device opens a WiFi AP with the `ssid:password` pair given as `microXT:a$$word`
2. Device listens as a TCP connection over the AP at `192.168.0.10` for two consecutive TCP requests (need to validate) on `port:9060`
3. Device connects to given WiFI provision
4. Device listens to TCP connection on `IP Address: 192.168.0.176:9060`
5. Device and client communicate via Google Protobuf