class Panel extends Element{
    constructor(override_id = ""){
        super("Panel", "panel", override_id)
        this.key ="Panel"
        this.updated = true;
    }

    mount(){
        super.mount();
    }

    html(){
        return `<div id=${this.id} class='panel'> <div class="medium config"><span class="bi-gear"></span></div>${stringify(this.children.map(item => item.html()))}</div>`
    }

}



function collectionCallback(event){
    //find the target component and update its children
    let component = mApp.getComponentById(event.currentTarget.id);
    console.log(`clicked: ${event.target.id}`)
    if(component !=null){
        $(`#${event.currentTarget.id} > span`).toggleClass("bi-chevron-right bi-chevron-down")
        component.children.forEach(item => $(`#${item.id}`).toggle());
    }else{
        console.log(`collection not found: ${event.target.id}`)
    }
}


class ItemCollection extends Element{
    constructor(key,name, hidden = false, callback = function(e){
        collectionCallback(e)
    }){
        super(key,name)
        this.count = 0;
        this.hidden = hidden;
        this.callback = callback
    }

     mount(){

        let binding = new Binding("click",this.callback);
        this.bindings = [binding]
        super.mount();

    }



    html(){
        var str = ""
        if(this.hidden){ str = "hidden"}
        return `<div id=${this.id} class='item small ${str}'><span class="small bi-chevron-right"></span><div class="subject">${this.name}</div><div class="indicator">${this.children.length}</div></div><div class="nest">${stringify(this.children.map(item => item.html()))}</div>`
    }

}


class Item extends Element{  
    constructor(key, name,hidden = false, callback = function(e){}){
        super(key,name)
        this.count = 0;
        this.id = name;
        this.callback = callback
        this.hidden = hidden;
    }

    mount(){
        
       let binding = new Binding( "click", this.callback)

        this.bindings.push(binding);
        super.mount();

    }




    html(){
        var str = ""
        if(this.hidden){ str = "hidden"}
        return `<div id=${this.id} class='item small ${str}'><span class="small bi-device-ssd"></span>${this.name}</div>`
    }

}
