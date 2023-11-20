//Parses json message for validation confirmation of wifi
//animates the successful call and proceeds to call a scaffolding
//requests
function endpointWifiValid(message){
    //Check JSON object for key valid
    if(message.extension.valid != null){
        if(message.extensions.valid == true){
            //Animate the wifi icon
            $("#wifi").removeClass("hidden")
            //Request the scaffolding
            astilectron.sendMessage({"name":"@scaffold"}, function(message){})
        }
    }
}


function endpointHasBluetooth(message){
    //Check JSON object for key valid
    if(message.extensions.hasBluetooth != null){
        if(message.extensions.hasBluetooth == "true"){
            //Animate the wifi icon
            $("#bluetooth-service").addClass("ico-info")
        }else{
            $("#bluetooth-service").removeClass("ico-info")
            $("#bluetooth-service").addClass("ico-std")
        }
    }
}


function endpointBluetoothScanning(message){
    //Check JSON object for key valid
    var interval = null
    if(message.extensions.scanning != null){
        if(message.scanning == "true"){
            interval  = pulseTextColor($("#bluetooth-scan"), "#2962b9", "#ffffff", 1000)
        }else{
            //Cancel the pulse color interval function
            clearInterval(interval)
        }
    }
}

function endpointBluetoothDisconnected(message){
    //Check JSON object for key valid
    if(message.extensions.disconnected != null){
        if(message.disconnected == "true"){
            //Animate the wifi icon
            $("#bluetooth-service").addClass("ico-err")
        }
    }
}

