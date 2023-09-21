/*App Hook -- Class is responsible for evaluating and updating component state and attaching
components to the DOM*/
class Sharp{
    constructor(){
        this.components = [];
        this.list_components = [];
        this.updated = true;
        this.attached = false;
        this.recieve = "";
        this.proxyEvents = new Map()
    }


    //setup the app by attaching component update hooks to their respective elements
    init(){

        for(let i = 0; i < this.components.length; i++){
            this.list_components = this.list_components.concat(this.components[i].getComponents())
            this.list_components.push(this.components[i])
        }

        $.when($.ready).then(function(){
            mApp.render()
            mApp.mount()
        })

    }

    mount(){
        for(let i = 0; i < this.list_components.length; i++){
            this.list_components[i].mount();
        }
    }

    add(component){
        this.components.push(component);
    }

    deleteBy(key){
        for(let i =0; i < this.components.length; i++){
            if (this.components[i].key === key){
                this.components.splice(i,1);
            }
        }
    }


    registerProxy(proxy_event){
        let key = `${proxy_event.event}:${proxy_event.proxy}`
        this.proxyEvents.set(key, proxy_event.component)
    }


    proxy(key){
        return this.proxyEvents.get(key)
    }

    deleteProxyEvent(key){
        this.proxyEvents.delete(key)
    }

    render(){
        for(let i = 0; i < this.components.length; i++){
            if(this.components[i].updated){
                let component = this.components[i];
                $("#"+component.id).html(component.html());
                component.updated = false;
            }
        }
    }


    getComponentById(id){
        for(let i = 0; i < this.components.length; i++){
            let item = this.components[i].getComponentById(id);
            if(item != null){
                return item;
            }
        }
        return null;
    }
}

class ProxyEvent{
    constructor(component, event_type, proxy_id){
        this.component = component
        this.event = event_type
        this.proxy = proxy_id
    }
}

class Binding{
    constructor(key, event){
        this.key = key
        this.event = event
    }
}

class KComponent{
    constructor(){
        this._name = "Component";
        this.key = "";
        this.class = "";
        this.id = "";
        this.updated = true;
        this.message = "";
        this.children = [];
        this.bindings = [];
        this.hidden = false;
    }

    mount(){
        for(let i = 0; i < this.bindings.length; i++){  
            let item = this.bindings[i];
            $(`#${this.id}`).on(item.key, item.event);
        }
    }d

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


var mApp = new Sharp(); //global - sharp
