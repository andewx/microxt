/*Sharp - Functional Javascript Frontend Endpoint Framework*/

/*DOM elements get added by - id=Name+UID, class=GroupName+GUID*/


class Session{
    constructor(){
        this.session = {}
    }

    set(key, value){
        this.session[key] = value
    }

    get(key){
        return this.session[key]
    }

    remove(key){
        delete this.session[key]
    }


    read(message){
        for(let key in message.session){
            console.log("SESSION:"+key + ":" + message.session[key])
            this.set(key, message.session[key])
        }
    }

    request(route, paramaters){
        var msg = {
            routekey:route,
            sessionkey:this.get("uid"),
            params:paramaters, //map[string]string
        
        }
        return JSON.stringify(msg)
    }
}

// Sharp is an endpoints and session frontend framework
// Exposes its enpoints with function handlers to JSON messages

class Sharp{
    constructor(sess){
        this.session = sess
        this._guid = 0;
        this.stopwatch = new Stopwatch(50/1000); //50ms frame timer
        this.endpoints = new Map();
    }

    guid(name){
        return name+this._guid++;
    }

    addEndpoint(name, handler){
        this.endpoints.set(name, handler)
    }

    call(name, message){
        this.endpoints.get(name)(message)
    }

    sessionRead(json){
        for(let key in json.session){
            this.session.set(key, json.session[key])
        }
    }


    //setup the app by attaching component update hooks to their respective elements
    init(){
        //sharp expects that the server will determine the initial state of the application

        global_session.set("uid", "0")
        $(document).ready(function(){
            htmx.defineExtension('hx-el', {
                onEvent : function(name,evt){
                  var msg = {
                    routekey:name,
                    sessionkey:global_session.get("uid"),
                    params:{}, //map[string]string
                  }

                  if(name === "@provision"){
                    msg.params["ssid"] = $("#ssid").val()
                    msg.params["password"] = $("#password").val()
                  }
                  astilectron.sendMessage(msg, function(message){
                    console.log(evt)
                  })
                }
              })
    
            //Add the root event listener for all components
            document.addEventListener("astilectron-ready", function(){
                astilectron.onMessage(function(message){
                    console.log(message)
                    if(message == null){
                        return
                    }
                    //Process all messages from GO - we expect JSON format
                    const json_message = JSON.parse(message)
                    if(json_message.type === "@endpoint"){
                        this.call(json_message.extensions.name, json_message)
                    }else if(json_message.type === "@error"){
                        console.log(json_message.extensions.error)
                    }else if(json_message.type === "@session"){
                        this.sessionRead(json_message)
                    }
                })

                //Make an initial session request to the server expecting @session response
                console.log(global_session.request("@session", {default:"default"}))
                astilectron.sendMessage(global_session.request("@session"), function(message){
                    //During response we request to get the initial view scaffold
                      astilectron.sendMessage(global_session.request("@scaffold"), function(message){})
                })

            })


            //Register the scaffolding by sending a request to the application

        })

    }

    event(callback){
        //global application event callback
    }

}

function stringify(obj){
    str =""
    for(let i = 0; i < obj.length; i++){
        str += obj[i].toString()
    }

    return str

}


class Stopwatch{
    constructor(period){
        this.period = period;
        this.start = new Date();
        this.end = null;
    }

    sample(){
        this.end = new Date();
    }

    getElapsed(){
        return this.end - this.start;
    }

    reset(){
        this.start = new Date();
        this.end = null;
    }

    hook(){
        this.sample();
        let elapsed = this.getElapsed();
        if(elapsed > this.period){
            this.reset();
            return true;
        }
        return false;
    }
}