package application

import (
	"encoding/json"
	"fmt"

	"github.com/andewx/microxt/app/templates"
	"github.com/asticode/go-astilectron"
)

type DomElement struct {
	HTML     string `json:"html"`
	Selector string `json:"selector"`
}

// Responds to a @scaffold request
func ScaffoldController(params map[string]string, session *Session, app Application) {
	req := NewRequest("@endpoint", session)
	html := &templates.StringWriter{Str: ""}
	tmpl := app.GetTemplate(session.state.ActiveView)
	var msg []byte
	var err error
	if tmpl != nil {
		tmpl.Execute(session.state, html)
	} else {
		fmt.Printf("Failed to find template %s\n", session.state.ActiveView)
	}

	req.Extensions["name"] = "@dom"
	req.Extensions["selectors"] = []*DomElement{{HTML: html.Str, Selector: "#root"}}

	//Send session request information to electron
	msg, err = json.Marshal(req)
	if err != nil {
		fmt.Printf("Failed to marshal session request %s\n", err.Error())

	} else {
		app.GetElectron().SendMessage(string(msg), func(m *astilectron.EventMessage) {})
	}

}

func ConnectDeviceController(params map[string]string, session *Session, app Application) {
	err := app.ConnectDevice(session.UID)
	req := NewRequest("@error", session)
	var msg []byte
	if err != nil {
		req.Extensions["error"] = "Failed to connect to device"
		msg, err = json.Marshal(req)
		if err != nil {
			app.GetElectron().SendMessage(string(msg), func(m *astilectron.EventMessage) {})
		} else {
			fmt.Printf("Failed to marshal error request %s\n", err.Error())
		}
	} else {
		session.state.ActiveView = "Ide"
		ScaffoldController(params, session, app)
	}
}

// Responds to a @provision request - for this we write the credentials over bluetooth connection
func ProvisionController(params map[string]string, session *Session, app Application) {
	ssid := params["ssid"]
	pass := params["pass"]
	var err error
	if ssid != "" && pass != "" {
		err = app.SendCredentials(ssid, pass, session) //Blocking call waits on a bluetooth connection
	}

	if err != nil {
		var msg []byte
		req := NewRequest("@error", session)
		req.Extensions["error"] = "Failed to connect to device"
		msg, err = json.Marshal(req)
		if err != nil {
			app.GetElectron().SendMessage(string(msg), func(m *astilectron.EventMessage) {})
		} else {
			fmt.Printf("Failed to marshal error request %s\n", err.Error())
		}
	} else {
		session.state.ActiveView = "Ide"
		ScaffoldController(params, session, app)
	}
}

// Responds to a @session request - this should establish a new session with the electron caller
func SessionController(params map[string]string, session *Session, app Application) {
	var msg []byte
	var err error
	req := NewRequest("@session", session)
	msg, err = json.Marshal(req)
	if err != nil {
		fmt.Printf("Failed to marshal session request %s\n", err.Error())

	} else {
		app.GetElectron().SendMessage(string(msg), func(m *astilectron.EventMessage) {})
	}

}

// Responds to a @devicepanel request
func DevicePanelController(params map[string]string, session *Session, app Application) {
	req := NewRequest("@dom", session)
	html := &templates.StringWriter{Str: ""}
	tmpl := app.GetTemplate("Devices")
	var msg []byte
	var err error
	if tmpl != nil {
		tmpl.Execute(session.state, html)
	} else {
		fmt.Printf("Failed to find template %s\n", session.state.ActiveView)
	}

	req.Extensions["selectors"] = []*DomElement{{HTML: html.Str, Selector: "#left-panel"}}

	//Send session request information to electron
	msg, err = json.Marshal(req)
	if err != nil {
		fmt.Printf("Failed to marshal session request %s\n", err.Error())

	} else {
		app.GetElectron().SendMessage(string(msg), func(m *astilectron.EventMessage) {})
	}

}

// Responds to a @navtabs request
func NavTabsController(params map[string]string, session *Session, app Application) {
	req := NewRequest("@dom", session)
	html := &templates.StringWriter{Str: ""}
	tmpl := app.GetTemplate("Nav")
	var msg []byte
	var err error
	if tmpl != nil {
		tmpl.Execute(session.state, html)
	} else {
		fmt.Printf("Failed to find template %s\n", session.state.ActiveView)
	}

	req.Extensions["selectors"] = []*DomElement{{HTML: html.Str, Selector: "#left-panel"}}

	//Send session request information to electron
	msg, err = json.Marshal(req)
	if err != nil {
		fmt.Printf("Failed to marshal session request %s\n", err.Error())

	} else {
		app.GetElectron().SendMessage(string(msg), func(m *astilectron.EventMessage) {})
	}

}

// Responds to a @terminaldisplay request
func TerminalDisplayController(params map[string]string, session *Session, app Application) {
	req := NewRequest("@dom", session)
	html := &templates.StringWriter{Str: ""}
	tmpl := app.GetTemplate("Terminal")
	var msg []byte
	var err error
	if tmpl != nil {
		tmpl.Execute(session.state, html)
	} else {
		fmt.Printf("Failed to find template %s\n", session.state.ActiveView)
	}

	req.Extensions["selectors"] = []*DomElement{{HTML: html.Str, Selector: "#left-panel"}}

	//Send session request information to electron
	msg, err = json.Marshal(req)
	if err != nil {
		fmt.Printf("Failed to marshal session request %s\n", err.Error())

	} else {
		app.GetElectron().SendMessage(string(msg), func(m *astilectron.EventMessage) {})
	}

}

// Responds to a @terminalinput request
func TerminalInputController(params map[string]string, session *Session, app Application) {
	/*
		req := NewRequest("@dom", session)
		//Look at terminal command and execute it
		command_string := params["command"]
		if command_string != ""{
			//Execute command
			cmd := NewCommand(command_string, session)
			cmd.Execute()
		}

		return
	*/
}

// Responds to a ADC Data Request
func RadarDataController(params map[string]string, session *Session, app Application) {

	/*var data []byte
	var err error
	req := NewRequest("@radardata", session)
	data_requested := params["type"]
	if data_requested == "adc" {
		data, err = json.Marshal(session.state.RadarData.ADCData)
	} else if data_requested == "fft" {
		data, err = json.Marshal(session.state.RadarData.FFTData)
	} else if data_requested == "pdat" {
		data, err = json.Marshal(session.state.RadarData.PDATData)
	} else if data_requested == "tdat" {
		data, err = json.Marshal(session.state.RadarData.TDATData)
	}

	if err != nil {
		msg := js.From(req).Begin("data").Set("data", data).End()
		app.GetElectron().SendMessage(msg, func(m *astilectron.EventMessage) {})
	} else {
		req := NewRequest("@error", session)
		tmp := js.From(req).Begin("error").Set("message", "Failed to marshal radar data").End()
		app.GetElectron().SendMessage(tmp.Marshal(), func(m *astilectron.EventMessage) {})
	}
	*/

}
