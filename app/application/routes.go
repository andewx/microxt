package application

type Route struct {
	Handler func(params map[string]string, session *Session, app *Application) error
}
