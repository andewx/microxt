package application

type Controller interface {
	Endpoint(name string, params map[string]string, session *Session, app *Application) error
}

type ControllerBase struct {
	handles map[string]func(params map[string]string, session *Session, app *Application) error
}

func NewControllerBase() *ControllerBase {
	return &ControllerBase{handles: make(map[string]func(params map[string]string, session *Session, app *Application) error)}
}

func (c *ControllerBase) Endpoint(name string, params map[string]string, session *Session, app *Application) error {

	if c.handles[name] != nil {
		return c.handles[name](params, session, app)
	}

	return nil
}
