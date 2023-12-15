
//Add the event handlers for the application
function handleProvisionRequest(evt){
    evt.preventDefault()
    const params = {"ssid":$("#ssid").val(), "pass":$("#pass").val()}
    console.log(sharp)
    const msg = sharp.request("BluetoothController","BluetoothProvisionWifi", params)
    astilectron.sendMessage(msg, function(message){
        console.log(evt)
    })
}

function handleProvisionCancel(evt){   
    evt.preventDefault()
    const msg = sharp.request("BluetoothController","Bluetooth Disconnect", params)
    astilectron.sendMessage(msg, function(message){
        console.log(evt)
    })

}


function handleScaffoldRequest(evt){ 
    evt.preventDefault()
    const msg = sharp.request("UtilityController","Scaffold", params)
    astilectron.sendMessage(msg, function(message){
        console.log(evt)
    })
}


//Handlers for the Register page
function handleUserSubmit(evt){
    evt.preventDefault()
    const params = {"first":$("#first").val(), "last":$("#last").val(), "email":$("#email").val(), "password":$("#password").val()}
    const msg = sharp.request("UserController","Register", params)
    astilectron.sendMessage(msg, function(message){
        console.log(evt)
    })
}

function handleUserCancel(evt){
    evt.preventDefault()
    const msg = sharp.request("UserController","RegisterCancel", params)
    astilectron.sendMessage(msg, function(message){
        console.log(evt)
    })
}

function preventDefault(evt){
    evt.preventDefault()
}