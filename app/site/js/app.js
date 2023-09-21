
//Application item callbacks
function configCallback(event){
    console.log("configCallback")
}

function schemaCallback(event){
    console.log("schemaCallback")
}

function profileCallback(event){
    console.log("profileCallback")
}

function validateCallback(event){
    console.log("validateCallback")
}

function analyzeCallback(event){
    console.log("analyzeCallback")
}

function connectivityCallback(event){
    console.log("connectivityCallback")
}

function addDeviceCallback(event){
    console.log("addDeviceCallback")
}


function frameworkCallback(event){
    console.log("frameworkCallback")
}


function captureCallback(event){
    console.log("captureCallback")
}

function playbackCallback(event){
    console.log("playbackCallback")
}


function openFilecallback(event){
    console.log("openFilecallback")
}


function saveFileCallback(event){
    console.log("saveFileCallback")
}


//Initial default left panel
group1 = new ItemCollection("Profile","Profiles")
profile1 = new ItemCollection("Devices","KLD7")
profile2 = new ItemCollection("Devices" ,"IWR6843")

device1 = new ItemCollection("Device","ArduinoMKR1010WiFi",true)
device2 = new ItemCollection("Device","InfineonK-LD724GhZ ranscieverDSPModule",true)

config1 = new Item("Action","Config", true, configCallback)
schema1 = new Item("Action","Schema", true,schemaCallback)
test1 = new Item("Action","Test",true, profileCallback)

config2 = new Item("Action","Config",true, configCallback)
schema2 = new Item("Action","Schema",true, schemaCallback)
test2 = new Item("Action","Test",true, profileCallback)

device1.children.push(config1)
device1.children.push(schema1)
device1.children.push(test1)

device2.children.push(config2)
device2.children.push(schema2)
device2.children.push(test2)

profile1.children.push(device1)
profile2.children.push(device2)

group1.children.push(profile1)
group1.children.push(profile2)


panel = new Panel("panel")
panel.children.push(group1)


//Fixed divider for our main application window
divider = new Divider("VerticalDrag", "divider", "divider")
mApp.add(divider)
mApp.add(panel)

mApp.init()

mApp.render()
