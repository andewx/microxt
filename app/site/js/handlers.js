
//Add the event handlers for the application
function handleProvisionRequest(evt){
    evt.preventDefault()
    const params = {"ssid":$("#ssid").val(), "password":$("#pass").val()}
    console.log(sharp)
    const msg = sharp.request("@provision", params)
    astilectron.sendMessage(msg, function(message){
        console.log(evt)
    })
}


function handleScaffoldRequest(evt){ 
    evt.preventDefault()
    const msg = sharp.request("@scaffold", params)
    astilectron.sendMessage(msg, function(message){
        console.log(evt)
    })
}