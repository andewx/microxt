class Divider extends Element{  
    constructor(key, name, override_id = ""){
        super(key,name, override_id)
    }

    mount(){
       var binding = new Binding("mousedown", startDrag);
        this.bindings = [binding]
        super.mount();
    }

    html(){
        return "" //for now this is a static existing element name of divider must match static index.html name
    }
}

/*Dragging functionality*/
function startDrag(e){
 var element = $("#main")
 var diff_y = e.pageY - element.height()
 var dragFunction = dragDivider(element, diff_y)
 var dropFunction = dropDivider(dragFunction)
 window.addEventListener("mousemove", dragFunction)
 window.addEventListener("mouseup", dropFunction)
 console.log("added listeners")
}


function dragDivider(component, diff_y){
    return function(e){
        h = e.pageY-diff_y
        component.height(h)
    }
}

function dropDivider(dragFunction){
    return function _drop(e){
        window.removeEventListener("mousemove", dragFunction)
        window.removeEventListener("mouseup", _drop)
    }
}
