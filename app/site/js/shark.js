/*App Hook -- Class is responsible for evaluating and updating component state and attaching
components to the DOM*/
class Shark{
    constructor(){
        this.components = [];
        this.list_components = [];
        this.updated = true;
        this.attached = false;
        this.recieve = "";
    }


    //setup the app by attaching component update hooks to their respective elements
    init(){

        for(let i = 0; i < this.components.length; i++){
            this.list_components = this.list_components.concat(this.components[i].getComponents())
            this.list_components.push(this.components[i])
        }

        $.when($.ready).then(function(){
            mApp.render()
            mApp.setup()
        })

    }

    setup(){
        for(let i = 0; i < this.list_components.length; i++){
            this.list_components[i].setup();
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

class Binding{
    constructor(){
        this.key = "";
        this.component = "";
        this.event = function(){}
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

    setup(){
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


var mApp = new Shark();