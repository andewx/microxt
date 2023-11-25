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



function endpointBluetoothOn(message){
    //Check JSON object for key valid
    $("#bluetooth-service").removeClass("ico-std")
    $("#bluetooth-service").addClass("ico-info")
}


function endpointBluetoothOff(message){
    //Check JSON object for key valid
    clearInterval(bluetoothIntervalID)
    $("#bluetooth-service").removeClass("ico-info")
    $("#bluetooth-service").addClass("ico-std")
    $("#bluetooth-scan").addClass("ico-std")
    $("#bluetooth-scan").css("color", "#ffffff")
    $("bluetooth-connect").removeClass("ico-info")
    $("bluetooth-connect").addClass("ico-std")
    $("bluetooth-disconnect").removeClass("ico-err")
    $("bluetooth-disconnect").addClass("ico-std")
}



function endpointBluetoothConnected(message){
    //Check JSON object for key valid
    if(message.extensions.connected != null){
        if(message.extensions.connected == "true"){
            //Animate the wifi icon
            $("#bluetooth-connect").addClass("ico-info")
        }else{
            $("#bluetooth-connect").removeClass("ico-info")
            $("#bluetooth-connect").addClass("ico-std")
        }
    }
}

var bluetoothIntervalID = null
function endpointBluetoothScanning(message){
    //Check JSON object for key valid
    var interval = null
    if(message.extensions.scanning != null){
        if(message.extensions.scanning == "true"){
            console.log("Scanning for bluetooth devices")
            const element = $("#bluetooth-scan")
            element.removeClass("ico-std")
            bluetoothIntervalID  = pulseTextColor(element, "#2962b9", "#ffffff", 50)
        }else{
            //Cancel the pulse color interval function
            clearInterval(bluetoothIntervalID)
        }
    }else{
        console.log("Message does not contain scanning key")
    }
}

function endpointBluetoothDisconnected(message){
    //Check JSON object for key valid
    if(message.extensions.disconnected != null){
        if(message.disconnected == "true"){
            //Animate the wifi icon
            $("#bluetooth-disconnect").addClass("ico-err")
        }
    }
}

function endpointDom(message){
    message.extensions.selectors.forEach((item)=>{
        //Set the inner html of the element
        $(item.selector).html(item.html)
    })
}

function endpointError(message){
    console.log(message)
}

function endpointRecieveFrame(message){

}

