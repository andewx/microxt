package application

import (
	"github.com/asticode/go-astilectron"
	js "github.com/coyove/jsonbuilder"
)

func NewScaffoldRoute(key string, mgr Routes) {
	route := Route{Key: key, Handler: nil}
	route.Handler = ScaffoldController
}

func ScaffoldController(params map[string]string, session *Session, app Application) {
	req := NewRequest("@update_dom", session)
	tmp := js.From(req)
	tmp.Begin("selectors").Begin("u0").Set("html", app.GetTemplate(session.State.ActiveView).Get()).Set("selector", "#root").End().End()
	app.GetElectron().SendMessage(tmp.Marshal(), func(m *astilectron.EventMessage) {
		//Response from m is ignored
	})
}
