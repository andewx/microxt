package application

import (
	"fmt"

	"github.com/andewx/microxt/app/templates"
	"github.com/asticode/go-astilectron"
)

type DomElement struct {
	HTML     string `json:"html"`
	Selector string `json:"selector"`
}

type UtilityController struct {
	*ControllerBase
}

func NewUtilityController() *UtilityController {
	c := &UtilityController{NewControllerBase()}
	c.handles["Scaffold"] = c.Scaffold
	c.handles["Session"] = c.Session
	c.handles["SendNotification"] = c.SendNotification
	c.handles["DevicePanel"] = c.DevicePanel
	c.handles["RadarData"] = c.RadarData
	return c
}

// Responds to a @scaffold request
func (u *UtilityController) Scaffold(params map[string]string, session *Session, app *Application) error {
	req := NewRequest("@endpoint", session)
	html := &templates.StringWriter{Str: ""}
	tmpl := app.GetTemplate(app.GetActiveView())

	if tmpl != nil {
		tmpl.Execute(session.state, html)
	} else {
		fmt.Printf("Failed to find template %s\n", app.GetActiveView())
	}

	req.Extensions["name"] = "@dom"
	req.Extensions["selectors"] = []*DomElement{{HTML: html.Str, Selector: "#root"}}

	//Send session request information to electron
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil

}

func (u *UtilityController) SendNotification(params map[string]string, session *Session, app *Application) error {
	req := NewRequest("@notification", session)
	req.Extensions["message"] = params["message"]
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil

}

// Responds to a @session request - this should establish a new session with the electron caller
func (u *UtilityController) Session(params map[string]string, session *Session, app *Application) error {
	req := NewRequest("@session", session)
	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil

}

// Responds to a @devicepanel request
func (u *UtilityController) DevicePanel(params map[string]string, session *Session, app *Application) error {
	req := NewRequest("@dom", session)
	html := &templates.StringWriter{Str: ""}
	tmpl := app.GetTemplate("DEVICE_VIEW")

	if tmpl != nil {
		tmpl.Execute(session.state, html)
	} else {
		fmt.Printf("Failed to find devices template\n")
	}

	req.Extensions["selectors"] = []*DomElement{{HTML: html.Str, Selector: "#left-panel"}}

	app.GetElectron().SendMessage(req.JSON(), func(m *astilectron.EventMessage) {})
	return nil

}

// Responds to a ADC Data Request
func (u *UtilityController) RadarData(params map[string]string, session *Session, app *Application) error {

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
	return nil

}
