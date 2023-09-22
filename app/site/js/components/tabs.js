class TabsBar extends Element{
    constructor(styles="tabs-bar",override_id = ""){
        super("TabsContainer", "tabscontainer", override_id)
        this.key ="TabsContainer"
        this.updated = true;
        this.styles = styles
    }

    mount(){
        super.mount();
    }

    html(){
        return `<div id=${this.id} class='tabs-bar'>${stringify(this.children.map(item => item.html()))}</div>`
    }

}



function tabSelect(event){
    //find the target component and update its children
    let component = mApp.getComponentById(event.currentTarget.id);
    if(component !=null){
        if(!$("#"+component.id).hasClass("active")){
            $("#"+component.id).toggleClass("active");
            component.parent.children.forEach(item =>{$(`#${item.id}`).toggleClass("active")});
        }
    }else{
        console.log(`Component not found: ${event.target.id}`)
    }
}


class Tab extends Element{
    constructor(key,name,parent, pageid, icon="bi-asterisk",active = false, callback = function(e){
        tabSelect(e)
    }){
        super(key,name)
        this.count = 0;
        this.callback = callback;
        this.icon = icon;
        this.active = active;
        this.pageid = pageid;
        this.parent = parent
    }

     mount(){

        let binding = new Binding("click",this.callback);
        this.bindings = [binding]
        super.mount();

    }



    html(){
        var str = ""
        var ind = ""
        if(this.active){ 
            str = "active"
            ind = "bi-snowflake"
        }
        return `<div id=${this.id} class='tab ${str}'><div class="subject">${this.name}</div><span class="indicator ${ind}"></span></div>`
    }

}
