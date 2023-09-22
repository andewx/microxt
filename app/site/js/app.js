
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


//Example Tree
group1 = new ItemCollection("Group","Devices")
profile1 = new ItemCollection("Profile","AirXT", "bi-diagram-2", true)
profile2 = new ItemCollection("Profile" ,"AirXT2","bi-diagram-2", true)

device1 = new ItemCollection("Device","AirXT Arduino MKR 1010","bi-cpu",true)
radar1 = new ItemCollection("Radar","Radar K-LD7 24GhZ DSP Module","bi-radar",true)

device2 = new ItemCollection("Device", "AirXT2 Arduino Portenta H7 DSP","bi-cpu",true)
radar2 = new ItemCollection("Radar", "DieselX mXT8 24GhZ ADC Module","bi-radar", true)

config1 = new Item("Action","Config", "bi-asterisk", "#898989ff", true, configCallback)
schema1 = new Item("Action","Schema", "bi-rainbow","#898989ff",true,schemaCallback)
test1 = new Item("Action","Test","bi-check-circle","green",true, profileCallback)

config2 = new Item("Action","Config", "bi-asterisk","#898989ff",true, configCallback)
schema2 = new Item("Action","Schema", "bi-rainbow","#898989ff",true,schemaCallback)
test2 = new Item("Action","Test","bi-check-circle","green",true, profileCallback)

config3 = new Item("Action","Config", "bi-asterisk","#898989ff",true, configCallback)
schema3 = new Item("Action","Schema", "bi-rainbow","#898989ff",true,schemaCallback)
test3 = new Item("Action","Test","bi-check-circle","green",true, profileCallback)

config4 = new Item("Action","Config", "bi-asterisk","#898989ff",true, configCallback)
schema4 = new Item("Action","Schema", "bi-rainbow","#898989ff",true,schemaCallback)
test4 = new Item("Action","Test","bi-check-circle","green",true, profileCallback)

device1.children.push(config1)
device1.children.push(schema1)
device1.children.push(test1)

radar1.children.push(config2)
radar1.children.push(schema2)
radar1.children.push(test2)

device2.children.push(config3)
device2.children.push(schema3)
device2.children.push(test3)

radar2.children.push(config4)
radar2.children.push(schema4)
radar2.children.push(test4)

profile1.children.push(device1)
profile1.children.push(radar1)

profile2.children.push(device2)
profile2.children.push(radar2)

group1.children.push(profile1)
group1.children.push(profile2)

panel = new Panel("panel")
panel.children.push(group1)


//Add some default tabs for display
tabsbar = new TabsBar("","maintabs")

tab1 = new Tab("tab", "Device Configuration", tabsbar,"page1","bi-egg", true)
tab2 = new Tab("tab", "MKR Wifi 101 Settings", tabsbar,"page1","bi-egg", false)
tab3 = new Tab("tab", "K-LD7 Settings", tabsbar,"page1","bi-egg", false) 
tab4 = new Tab("tab", "Network Settings", tabsbar,"page1","bi-egg", false) 
tab5 = new Tab("tab", "Render", tabsbar,"page1","bi-egg", false) 

tabsbar.children.push(tab1,tab2,tab3,tab4,tab5)

//Fixed divider for our main application window
divider = new Divider("VerticalDrag", "divider", "divider")
mApp.add(divider)
mApp.add(panel)
mApp.add(tabsbar)

mApp.init()

mApp.render()
