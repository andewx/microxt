package application

import (
	"math/rand"
	"strconv"

	"github.com/andewx/microxt/app/models"
	"github.com/andewx/microxt/app/templates"
	"github.com/andewx/microxt/common"
	"github.com/andewx/microxt/net"
	"github.com/asticode/go-astilectron"
)

// For now we will keep the Application interface as simple as possible. To facilitate ease of usage
// Conceptually the Application interface manages a set of templates, views, routes, and sessions.
// When an application processes a route handler should have resources to access the session and the
// tcp/electron interfaces. The application event loop should be executing endpoints for both DeviceListeningTCP Ports
// and executing endpints for AstiElectron event messages
type Application interface {
	AddTemplate(key string, temp *templates.ApplicationTemplate) error
	AddRoute(key string, rte *Route) error
	GetNet() *net.TCPConnection
	GetElectron() *astilectron.Window
	GetSession(key string) *Session
	GetTemplate(key string) *templates.ApplicationTemplate
	GetRoute(key string) *Route
	GetKeys() Keys
	SetWindow(*astilectron.Window)
	ConnectDevice(string) error
	SendCredentials(string, string, *Session) error
	GetChannel() chan int
	NewSession() *Session
	Close()
}

type Request struct {
	Type       string                 `json:"type"`
	Session    *Session               `json:"session"`
	Extensions map[string]interface{} `json:"extensions"`
}

type Session struct {
	UID   string `json:"uid"`
	state *models.SessionObject
}

type Keys map[string]*Session

func NewSession(view string, mgr Keys) *Session {
	x := int(rand.Int63())
	strx := strconv.FormatInt(int64(x), 16)
	sesh := &Session{UID: strx, state: models.NewSessionObject()}
	if mgr[strx] == nil {
		mgr[strx] = sesh
	}
	return sesh
}

func NewRequest(_type string, sesh *Session) *Request {
	return &Request{Type: _type, Session: sesh, Extensions: make(map[string]interface{}, 0)}
}

// Air Application
type AirApplication struct {
	AppConfig     *models.AppConfig
	Sessions      Keys
	Routes        map[string]*Route
	TemplateViews map[string]*templates.ApplicationTemplate
	Connection    *net.TCPConnection
	Window        *astilectron.Window
	_state        int
	DeviceChannel chan int
}

func NewAirXTApplication() (*AirApplication, error) {
	var app = new(AirApplication)
	var err error
	app.AppConfig = models.NewAppConfig()
	app.Sessions = make(Keys)
	app.Routes = make(map[string]*Route, 0)
	app.TemplateViews = make(map[string]*templates.ApplicationTemplate, 0)
	app.DeviceChannel = make(chan int)

	//AddRoutes
	app.AddRoute("@provision", &Route{Handler: ProvisionController})
	app.AddRoute("@session", &Route{Handler: SessionController})
	app.AddRoute("@devicepanel", &Route{Handler: DevicePanelController})
	app.AddRoute("@navtabs", &Route{Handler: NavTabsController})
	app.AddRoute("@terminaldisplay", &Route{Handler: TerminalDisplayController})
	app.AddRoute("@terminalinput", &Route{Handler: TerminalInputController})
	app.AddRoute("@radardata", &Route{Handler: RadarDataController})
	app.AddRoute("@scaffold", &Route{Handler: ScaffoldController})

	//Add Templates
	app.AddTemplate("Provision", templates.NewTemplate("Provision", common.ProjectRelativePath("microxt/app/templates/provision.gohtml")))
	app.AddTemplate("Login", templates.NewTemplate("Login", common.ProjectRelativePath("microxt/app/templates/login.gohtml")))
	app.AddTemplate("Ide", templates.NewTemplate("Ide", common.ProjectRelativePath("microxt/app/templates/ide.gohtml")))

	app._state = net.NO_DEVICE
	str := app.AppConfig.ProvisionIP + ":" + strconv.FormatInt(int64(app.AppConfig.ProvisionPort), 10)
	app.Connection, err = net.NewTCPConnection(str) //If there is no TCP connection then there is no device we can try again later
	if err != nil {
		err = nil
	}
	return app, err
}

func (p *AirApplication) NewSession() *Session {
	//generate unique session id
	session := NewSession("Provision", p.GetKeys())
	return session
}

func (p *AirApplication) Close() {
	p.Connection.Close()
	p.DeviceChannel <- EXIT
}

func (p *AirApplication) GetChannel() chan int {
	return p.DeviceChannel
}

func (p *AirApplication) SetWindow(window *astilectron.Window) {
	p.Window = window
}

func (p *AirApplication) SendCredentials(ssid string, password string, session *Session) error {
	var err error
	//Send SSID/Password as raw byte messages to the connection in succession
	device, err := net.NewBluetoothConnection()
	if err != nil {
		done := false
		for !done {
			msg := <-device.Status
			if msg == net.DEVICE_CONNECTED {
				err = device.Write(net.SSID_CHARACTERISTIC, []byte(ssid))
				err = device.Write(net.PASS_CHARACTERISTIC, []byte(password))
				req := NewRequest("@endpoint", session)
				req.Extensions["name"] = "@bluetoothConnected"
				req.Extensions["connected"] = "true"
				p.GetElectron().SendMessage("@endpoint", func(m *astilectron.EventMessage) {})
				done = true
			} else if msg == net.DEVICE_DISCONNECTED {
				//Send endpoint message to the application
				req := NewRequest("@endpoint", session)
				req.Extensions["name"] = "@bluetoothDisconnected"
				req.Extensions["disconnected"] = "true"
				p.GetElectron().SendMessage("@endpoint", func(m *astilectron.EventMessage) {})

			} else if msg == net.DEVICE_SCANNING {
				req := NewRequest("@endpoint", session)
				req.Extensions["name"] = "@bluetoothScanning"
				req.Extensions["scanning"] = "true"
				p.GetElectron().SendMessage("@endpoint", func(m *astilectron.EventMessage) {})

			} else if msg == net.DEVICE_ON {
				req := NewRequest("@endpoint", session)
				req.Extensions["name"] = "@bluetoothOn"
				req.Extensions["valid"] = "true"
				p.GetElectron().SendMessage("@endpoint", func(m *astilectron.EventMessage) {})

			} else if msg == net.NO_DEVICE {
				req := NewRequest("@endpoint", session)
				req.Extensions["name"] = "@bluetoothOn"
				req.Extensions["valid"] = "false"
				p.GetElectron().SendMessage("@endpoint", func(m *astilectron.EventMessage) {})
			}
		}
	}
	device.Close()
	return err
}

func (p *AirApplication) ConnectDevice(session_id string) error {
	var err error
	p.Connection, err = net.NewTCPConnection(p.AppConfig.DeviceIP + ":" + strconv.FormatInt(int64(p.AppConfig.DevicePort), 10))
	if err != nil {
		p._state = net.DEVICE_CONNECTED
		session := p.GetSession(session_id)
		session.state.Devices = append(session.state.Devices, &models.Device{IP: p.AppConfig.DeviceIP, Port: p.AppConfig.DevicePort, Driver: "KLD7", Name: "24Ghz K-LD7 Radar"})
		go p.Connection.Listen(p.DeviceChannel)
	}
	return err
}

func (p *AirApplication) CloseDevice() {
	p.Connection.Close()
}

func (p *AirApplication) GetKeys() Keys {
	return p.Sessions
}

func (p *AirApplication) AddRoute(key string, rte *Route) error {
	p.Routes[key] = rte
	return nil
}

func (p *AirApplication) AddTemplate(key string, tmpl *templates.ApplicationTemplate) error {
	p.TemplateViews[key] = tmpl
	return nil
}

func (p *AirApplication) GetNet() *net.TCPConnection {
	return p.Connection
}

func (p *AirApplication) GetElectron() *astilectron.Window {
	return p.Window
}

func (p *AirApplication) GetSession(key string) *Session {
	return p.Sessions[key]
}

func (p *AirApplication) GetTemplate(key string) *templates.ApplicationTemplate {
	return p.TemplateViews[key]
}

func (p *AirApplication) GetRoute(key string) *Route {
	return p.Routes[key]
}
