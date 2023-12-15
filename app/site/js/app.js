
//Defines API Shape for JS Enpoint requests to TCP JSON Caller
//JSON Message Formats Ensure Extensible Core Routing Protocols can be adapted
sharp.addEndpoint("@wifi", endpointWifiValid)
sharp.addEndpoint("@bluetoothOn", endpointBluetoothOn)
sharp.addEndpoint("@bluetoothScanning", endpointBluetoothScanning)
sharp.addEndpoint("@bluetoothDisconnected", endpointBluetoothDisconnected)
sharp.addEndpoint("@bluetoothConnected", endpointBluetoothConnected)
sharp.addEndpoint("@bluetoothOff", endpointBluetoothOff)
sharp.addEndpoint("@bluetoothSuccess", endpointBluetoothSuccess)
sharp.addEndpoint("@dom", endpointDom)
sharp.addEndpoint("@error", endpointError)
sharp.addEndpoint("@recieveFrame", endpointRecieveFrame)
sharp.addEndpoint("@bluetoothCancel", endpointBluetoothOff)
bluejs.addBinding("provision", {}, handleProvisionRequest)
bluejs.addBinding("scaffold",{}, handleScaffoldRequest)
bluejs.addBinding("provisionCancel",{}, handleProvisionCancel)
bluejs.addBinding("userSubmit",{}, handleUserSubmit)
bluejs.addBinding("userCancel",{}, handleUserCancel)



$(document).ready(function(){
    document.addEventListener("astilectron-ready", function(){
        astilectron.onMessage(function(message){
            console.log(message)
            if(message == null){
                return
            }
            //Process all messages from GO - we expect JSON format
            const json_message = JSON.parse(message)
            if(json_message == null){
                console.log("Message could not be resolved to JSON" +message)
                return
            }
            if(json_message.type === "@endpoint"){
                sharp.call(json_message.extensions.name, json_message)
            }else if(json_message.type === "@error"){
                console.log(json_message.extensions.error)
            }else if(json_message.type === "@session"){
                sharp.sessionRead(json_message)
            }
        })
    })
    //Register the scaffolding by sending a request to the sharp

})