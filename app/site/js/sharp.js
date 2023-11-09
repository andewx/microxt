/*Sharp - Simple virtual component management, expects to reflect state from server and acts a JS attachment
for components retrieved from the server. Mount components attaches callbacks to the components. Render is
a utility method for user defined component effects given the Sharp Object-Event framerate loop update*/

/*DOM elements get added by - id=Name+UID, class=GroupName+GUID*/
class Sharp{
    constructor(){
        this.components = new Map()
        this.flat_components = new Map()
        this._guid = 0;
        this.stopwatch = new Stopwatch(50/1000); //50ms frame timer
    }

    guid(name){
        return name+this._guid++;
    }

    //setup the app by attaching component update hooks to their respective elements
    init(){
        //sharp expects that the server will determine the initial state of the application
        $(document).ready(function(){
            document.addEventListener("astilectron-ready", function(){
                astilectron.onMessage(function(message){
                    //Process all messages from GO - we expect JSON format
                    const json_message = JSON.parse(message)
                    if(json_message.type === "@update_dom"){
                        json_message.selectors.forEach((item)=>{
                            $(item.selector).html(item.html)
                        })
                      
                    }else if(json_message.type === "@refresh_components"){
                        this.clear()
                        //May have to sort the components first
                        if(json_message.variable === "components"){
                            json_message.value.forEach((element)=>{
                                newComponentFromJSON(element)
                            })
                        }
                    }else if(json_message.type === "@add_component"){
                       let element = json_message.value.element
                       newComponentFromJSON(element)
                    }else if(json_message.type === "@delete_component"){
                        let key = json_message.value
                        this.components.delete(key)
                    }else if(json_message.type === "@update_component"){
                        let key = json_message.key
                        let props = json_message.props
                        let element = json_message.element
                        let component = this.components.get(key)
                        component.update(element, props)
                    }
                })
            })
        })

        $.when($.ready).then(function(){
            mApp.mount()
        })

    }

    event(callback){
        //global application event callback
    }


    mount(){
        this.components.forEach((value, key, map) => {
            value.mount()
        })
    }


    get(key){
        return this.flat_components.get(key)
    }

    addByKey(component,key){
        this.components.set(key, component)
    }

    delete(key){
        this.flat_components.delete(key)
        this.components.delete(key)
    }

    render(){
        this.components.forEach((value, key, map) => {
            value.render()
        })
    }
}

function newComponentFromJSON(element){
    let component = null
    if(element.type === "Collection"){
        component = new ItemCollection(element.key, element.name, element.icon, element.hidden)
        this.addByKey(component, element.key)
    }
    else if(element.type === "Item"){
        component = new Item(element.key, element.name, element.icon, element.hidden)
        if(element.belongs === ""){
            this.addByKey(component, element.key)
        }
        belongs = this.map_components.get(element.belongs)
        belongs.children.push(component)
    }
    else if(element.type === "Tab"){
        component = new Tab(element.key, element.name, element.icon, element.hidden)
        if(element.belongs === ""){
            this.addByKey(component, element.key)
        }
        belongs = this.map_components.get(element.belongs)
        belongs.children.push(component)
    }
    else if(element.type === "Panel"){
        component = new Panel(element.key, element.name, element.icon, element.hidden)
        this.addByKey(component, element.key)
    }
}

class Binding{
    constructor(key, event){
        this.key = key
        this.event = event
    }
}

class Element{
    constructor(key, name, override_id = ""){
        this._name = key;
        this.name = name;
        if(override_id === ""){
            var name_nospace = name
            name_nospace = name_nospace.replace(/\s/g, '')
            this.id = mApp.guid(name_nospace);
        }else{
            this.id = override_id.replace(/\s/g,'');
        }

        if(this.id===""){
            console.log("id is empty fail!!!!")
        }
        this.updated = true;
        this.children = [];
        this.bindings = [];
        this.hidden = false;
    }


    update(element, props){

    }

    mount(){
        for(let i = 0; i < this.bindings.length; i++){  
            let item = this.bindings[i];
            $(`#${this.id}`).on(item.key, item.event);
        }

        for (let i = 0; i < this.children.length; i++){
            let component = this.children[i]
            component.mount()
        }
    }

    getComponents(){
            let arr = this.children.map(item => item)
            for(let i = 0; i < arr.length; i++){
                arr = arr.concat(arr[i].getComponents())
            }
        
            return arr
    }

    getComponentById(id){
        if(this.id === id){
            return this;
        }
        else{
            for(let i = 0; i < this.children.length; i++){
                let item = this.children[i].getComponentById(id);
                if(item != null){
                    return item;
                }
            }
        }
        return null;
    }

    addUnique(id){
        this.children.push(id);
    }

    addClass(className){
        this.class = className;
    }

    addBinding(key, binding){
        let newBinding = new Binding();
        newBinding.key = key;
        newBinding.event = binding;
        this.bindings.push(newBinding);
    }

    html(){
        var str =  ""
        for(let i = 0; i < this.children.length; i++){
            str += this.children[i].html();
        }
        return str
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


var mApp = new Sharp(); //global - sharp
