var global_session = new Session()
var application = new Sharp(global_session)
application.init()

application.addEndpoint("wifi", application.wifi)
application.addEndpoint("$bluetoothOn", application.bluetooth)
application.addEndpoint("$bluetoothScanning", application.bluetoothScan)
application.addEndpoint("$bluetoothDisconnected", application.bluetoothDisconnected)
application.addEndpoint("$bluetoothConnected", application.bluetoothConnected)
