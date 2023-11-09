class Panel extends Element{
    constructor(override_id = ""){
        super("Panel", "panel", override_id)
        this.key ="Panel"
        this.updated = true;
    }

    mount(){
        super.mount();
    }

    render(){
        super.render();
    }

    html(){
        return `<div id=${this.id} class='panel'> <div class="large config"><span class="bi-plus-circle"></span></div>${stringify(this.children.map(item => item.html()))}</div>`
    }

}



function collectionCallback(event){
    //find the target component and update its children
    let component = mApp.getComponentById(event.currentTarget.id);
    if(component !=null){
        console.log(`Component Toggle: ${event.currentTarget.id}`)
        $(`#${event.currentTarget.id} > span`).toggleClass("bi-chevron-right bi-chevron-down")
        component.children.forEach(item =>{$(`#${item.id}`).toggleClass("hidden")});
      
    }else{
        console.log(`Component not found: ${event.target.id}`)
    }
}


class ItemCollection extends Element{
    constructor(key,name, icon = "bi-folder",hidden = false, callback = function(e){
        collectionCallback(e)
    }){
        super(key,name)
        this.count = 0;
        this.hidden = hidden;
        this.callback = callback;
        this.icon = icon;
    }

     mount(){

        let binding = new Binding("click",this.callback);
        this.bindings = [binding]
        super.mount();

        
    }

    render(){
        super.render();
    }



    html(){
        var str = ""
        if(this.hidden){ str = "hidden"}
        return `<div id=${this.id} class='item ${str}'><span class="small bi-chevron-right"></span><span class="small ${this.icon}"></span><div class="subject">${this.name}</div><div class="indicator">${this.children.length}</div></div><div class="nest">${stringify(this.children.map(item => item.html()))}</div>`
    }

}


class Item extends Element{  
    constructor(key, name, icon ="bi-asterisk", hovercolor="white", hidden = false, callback = function(e){}){
        super(key,name)
        this.count = 0;
        this.callback = callback
        this.hidden = hidden;
        this.icon = icon
        this.color = hovercolor;
    }

    mount(){
        
       let binding = new Binding( "click", this.callback)

        this.bindings = [binding]
        super.mount();

    }


    render(){
        super.render();
    }

    html(){
        var str = ""
        if(this.hidden){ str = "hidden"}
        return `<div id=${this.id} class='item small ${str}' style="color:${this.color}"><span class="${this.icon}"  style="color:${this.color}"></span>${this.name}</div>`
    }

}
