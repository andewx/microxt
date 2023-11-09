class Terminal extends Element{
    constructor(styles="terminal",override_id = ""){
        super("terminalContainer", "terminalcontainer", override_id)
        this.key ="terminalContainer"
        this.updated = true;
        this.styles = styles
    }

    mount(){
        super.mount();
    }

    render(){
        super.render();
    }

    html(){
        return `<div id=${this.id} class='terminal'> <div class="large config"><span class="bi-plus-circle"></span></div>${stringify(this.children.map(item => item.html()))}</div>`
    }

}

class TerminalInput extends Element{
    constructor(styles="terminalinput",override_id = ""){
        super("TerminalInput", "terminalinput", override_id)
        this.key ="terminalinput"
        this.updated = true;
        this.styles = styles
    }

    mount(){
        super.mount();
    }

    render(){
        super.render();
    }

    html(){
        return `<div id=${this.id} class='terminal'> <div class="large config"><span class="bi-plus-circle"></span></div>${stringify(this.children.map(item => item.html()))}</div>`
    }

}
