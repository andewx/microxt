package application

type Routes map[string]*Route

type Route struct {
	Key     string
	Handler func(params map[string]string, session *Session, app Application)
}
