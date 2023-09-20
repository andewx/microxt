class FileView extends KComponent{
    constructor(){
        super()
        this.key ="FileView"
        this.id = "file-view";
        this.updated = true;
    }

}


class Folder extends KComponent{
    constructor(count, name){
        super()
        this.key ="Folder"
        this.count = count;
        this.name = name;
        this.id = "folder"+this.count;
        this.bindings = [{
            key: "click",
            function(){
            }
        }];
    }

     setup(){

       let binding = new Binding();
        binding.key = "click";
        binding.id = this.id;
        binding.event = function(event){
            //find the target component and update its children
            let component = mApp.getComponentById(event.currentTarget.id);
            console.log(`clicked: ${event.target.id}`)
            if(component !=null){
                $(`#${event.currentTarget.id} > span`).toggleClass("bi-chevron-right bi-chevron-down")
                component.children.forEach(item => $(`#${item.id}`).toggle());
            }else{
                console.log(`component not found: ${event.target.id}`)
            }
        }

        this.bindings = [binding]
        super.setup();

    }



    html(){
        return `<div id=${this.id} class='item small'><span class="small bi-chevron-right"></span><div class="subject">${this.name}</div><div class="indicator">${this.children.length}</div></div><div class="nest">${stringify(this.children.map(item => item.html()))}</div>`
    }

}


class File extends KComponent{  
    constructor(count, name){
        super()
        this.key ="File"
        this.count = count;
        this.name = name;
        this.id = "file"+this.count;
        this.hidden = true;
        this.updated = true;
    }

    setup(){
        
       let binding = new Binding();
        binding.key = "click";
        binding.event = function(event){
            console.log(`clicked: ${event.target.id}`)
        }
        this.bindings.push(binding);

        super.setup();

    }




    html(){
        return `<div id=${this.id} class='item small' style="display:none"><span class="small bi-device-ssd"></span>${this.name}</div>`
    }

}