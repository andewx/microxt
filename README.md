# README.md

### Welcome

**Overview**

Microxt is a multiplatform application developed upon Golang and the Electron Framework API for interfacing with IoT devices. The application scans for devices on the local area network operating on UDP out of port 8040 and connects to devices that respond properly to a handshake request which requires an public/private key match.

Microxt handles the request/response routing via key based messages, which specify content type or response sequences to follow. 