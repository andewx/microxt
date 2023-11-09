package application

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/andewx/microxt/app/models"
	"github.com/andewx/microxt/app/templates"
	"github.com/andewx/microxt/net"
	"github.com/asticode/go-astilectron"
)

// For now we will keep the Application interface as simple as possible. To facilitate ease of usage
// Conceptually the Application interface manages a set of templates, views, routes, and sessions.
// When an application processes a route handler should have resources to access the session and the
// tcp/electron interfaces. The application event loop should be executing endpoints for both DeviceListeningTCP Ports
// and executing endpints for AstiElectron event messages
type Application interface {
	Create(*astilectron.Window) error
	AddTemplate(key string, temp templates.Template) error
	AddRoute(rte *Route) error
	GetNet() *net.TCPConnection
	GetElectron() *astilectron.Window
	GetSession(key int64) *Session
	GetSessionByUserID(key int64) *Session
	GetTemplate(key string) templates.Template
	GetRoute(key string) *Route
}

type Request struct {
	_Type     string
	Timestamp int64
	Session   *Session
}

type Session struct {
	UID    int64
	UserID int64
	State  *models.SessionObject
}

type Keys map[int64]*Session

func NewSession(UserID int64, view string, mgr Keys) *Session {
	x := rand.Int63()
	sesh := &Session{UID: x, UserID: UserID, State: models.NewSessionObject()}
	for mgr[x] != nil {
		x = rand.Int63()
		if mgr[x] == nil {
			sesh.UID = x
			mgr[x] = sesh
		}
	}
	return sesh
}

func NewRequest(_type string, sesh *Session) *Request {
	return &Request{_Type: _type, Timestamp: time.Now().Unix(), Session: sesh}
}

// Air Application
type AirApplication struct {
	AppConfig     *models.AppConfig
	Sessions      Keys
	Routes        Routes
	TemplateViews map[string]templates.Template
	Connection    *net.TCPConnection
	Window        *astilectron.Window
}

func NewAirXTApplication() (*AirApplication, error) {
	var app = new(AirApplication)
	var err error
	app.AppConfig = models.NewAppConfig()
	app.Sessions = make(Keys)
	app.Routes = make(Routes)
	app.TemplateViews = make(map[string]templates.Template)
	app.Connection = net.NewTCPConnection(app.AppConfig.ProvisionIP + ":" + strconv.FormatInt(int64(app.AppConfig.ProvisionPort), 10))
	return app, err
}

func (p *AirApplication) Create(window *astilectron.Window) error {
	p.Window = window
	return nil
}

func (p *AirApplication) AddRoute(rte *Route) error {
	p.Routes[rte.Key] = rte
	return nil
}

func (p *AirApplication) AddTemplate(key string, tmpl templates.Template) error {
	p.TemplateViews[key] = tmpl
	return nil
}

func (p *AirApplication) GetNet() *net.TCPConnection {
	return p.Connection
}

func (p *AirApplication) GetElectron() *astilectron.Window {
	return p.Window
}

func (p *AirApplication) GetSession(key int64) *Session {
	return p.Sessions[key]
}

func (p *AirApplication) GetSessionByUserID(key int64) *Session {
	for _, sesh := range p.Sessions {
		if sesh.UserID == key {
			return sesh
		}
	}
	return nil
}

func (p *AirApplication) GetTemplate(key string) templates.Template {
	return p.TemplateViews[key]
}

func (p *AirApplication) GetRoute(key string) *Route {
	return p.Routes[key]
}
