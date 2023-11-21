var global_session = new Session()
var application = new Sharp(global_session)
application.init()

//Defines API Shape for JS Enpoint requests to TCP JSON Caller
//JSON Message Formats Ensure Extensible Core Routing Protocols can be adapted
application.addEndpoint("@wifi", endpointWifiValid)
application.addEndpoint("@bluetoothOn", endpointHasBluetooth)
application.addEndpoint("@bluetoothScanning", endpointBluetoothScanning)
application.addEndpoint("@bluetoothDisconnected", endpointBluetoothDisconnected)
application.addEndpoint("@bluetoothConnected", endpointHasBluetooth)
application.addEndpoint("@dom", endpointDom)
application.addEndpoint("@error", endpointError)
application.addEndpoint("@recieveFrame", endpointRecieveFrame)

