
class Divider extends KComponent{  
    constructor(count, name){
        super()
        this.key ="Divider"
        this.count = count;
        this.name = name;
        this.id = name+this.count;
        this.hidden = true;
        this.updated = true;
        this.currentY = 0;

    }

    setup(){
        
       let binding = new Binding();
        binding.key = "mousedown";
        binding.event = function(event){
            //bind dynamic handlers to mouse move and mouse up with the window
            //this will allow the user to drag the divider and resize the panes
           let component = mApp.getComponentById(event.target.id);
            component.currentY = event.pageY;

                let resizeWindow = function(event){
                    let h = component.currentY - event.pageY;
                    let y = $(`#main`).height()
                    component.currentY = event.pageY;
                    console.log(`y: ${y} h: ${h}`)
                    $(`#main`).height(y-h);
                }

                $(document.body).on("mousemove", resizeWindow(event));

                $(document.body).on("mouseup", function(event){
                    console.log("mouseup")
                    $(document.body).off("mousemove");
                    component.setup()
                })  
            
   
        }

        this.bindings = [binding]

        super.setup();

    }

    html(){
        return "" //for now this is a static existing element name of divider must match static index.html name
    }

}