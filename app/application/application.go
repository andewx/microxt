package application

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/andewx/microxt/app/models"
	"github.com/andewx/microxt/app/templates"
	"github.com/andewx/microxt/common"
	"github.com/andewx/microxt/net"
	"github.com/asticode/go-astilectron"
)

// View Enumerations
const (
	MAIN_VIEW = iota
	PROVISION_VIEW
	DEVICE_DETAILS_VIEW
	REGISTER_VIEW
	LOG_VIEW
	SUCCESS_VIEW
	LOAD_VIEW
	FATAL_VIEW
)

type Request struct {
	Type       string                 `json:"type"`
	Session    *Session               `json:"session"`
	Extensions map[string]interface{} `json:"extensions"`
}

type Session struct {
	UID   string `json:"uid"`
	state *models.SessionObject
}

func NewSession(view string, mgr map[string]*Session) *Session {
	x := int(rand.Int63())
	strx := strconv.FormatInt(int64(x), 16)
	sesh := &Session{UID: strx, state: models.NewSessionObject()}
	sesh.state.User = &models.User{}
	if mgr[strx] == nil {
		mgr[strx] = sesh
	}
	return sesh
}

func NewRequest(_type string, sesh *Session) *Request {
	return &Request{Type: _type, Session: sesh, Extensions: make(map[string]interface{}, 0)}
}

func (r *Request) JSON() string {
	var msg []byte
	var err error
	msg, err = json.Marshal(r)
	if err != nil {
		fmt.Printf("Failed to marshal session request %s\n", err.Error())
	}
	return string(msg)
}

/*
Application component manages Application resources and state for now we only allow one active

	connection however this may be extended to host multiple concurrent devices in the future if needed
*/
type Application struct {
	AppConfig        *models.AppConfig
	Sessions         map[string]*Session
	Controllers      map[string]Controller
	TemplateViews    map[string]*templates.ApplicationTemplate
	Connection       *net.TCPConnection
	Window           *astilectron.Window
	ActiveView       int
	Devices          map[string]*models.Device `json:"devices"`
	ConnectedDevices []*models.Device
	ActiveDevice     *models.Device
	BiChan           chan int
	User             *models.User `json:"user"`
}

func NewApplication() (*Application, error) {
	var app = new(Application)
	var err error
	app.AppConfig = models.NewAppConfig()
	app.Sessions = make(map[string]*Session)
	app.Controllers = make(map[string]Controller)
	app.TemplateViews = make(map[string]*templates.ApplicationTemplate, 0)
	app.Devices = make(map[string]*models.Device, 0)
	app.BiChan = make(chan int)

	//Add Controller
	app.AddController("UtilityController", NewUtilityController())
	app.AddController("MainController", NewMainController())
	app.AddController("UserController", NewUserController())
	app.AddController("BluetoothController", NewBluetoothController())
	app.AddController("DeviceController", NewDeviceController())

	//Add Templates
	if err = app.AddTemplate("PROVISION_VIEW", templates.NewTemplate("PROVISION_VIEW", common.ProjectRelativePath("microxt/app/templates/provision.tmpl"))); err != nil {
		return nil, err
	}

	if err = app.AddTemplate("REGISTER_VIEW", templates.NewTemplate("REGISTER_VIEW", common.ProjectRelativePath("microxt/app/templates/register.tmpl"))); err != nil {
		return nil, err
	}

	if err = app.AddTemplate("MAIN_VIEW", templates.NewTemplate("MAIN_VIEW", common.ProjectRelativePath("microxt/app/templates/main.tmpl"))); err != nil {
		return nil, err
	}

	app.ActiveView = MAIN_VIEW

	return app, err
}

func (p *Application) Init() error {
	var err error

	//Initiate a new session and send to the Electron
	session := p.NewSession()
	request := NewRequest("@session", session)
	err = p.Window.SendMessage(request.JSON(), func(m *astilectron.EventMessage) {})
	if err != nil {
		return err
	}

	//Attempt to load devices from the user.json file
	main := p.Controllers["MainController"].(*MainController)
	user := p.Controllers["UserController"].(*UserController)
	err = main.ReadDevices(nil, session, p)
	if err != nil || len(p.Devices) == 0 {
		fmt.Printf("Failed to load any devices\n")
		p.ActiveView = PROVISION_VIEW
	}

	//Attempt a connection to the most recently used device
	err = main.SelectMostRecentDevice(nil, session, p)
	if err != nil {
		fmt.Printf("Failed to select most recent device\n")
		p.ActiveView = PROVISION_VIEW
	}

	//Attempt to connect to the device
	err = main.ConnectActiveDevice(nil, session, p)
	if err != nil {
		fmt.Printf("Failed to connect to device\n")
		p.ActiveView = PROVISION_VIEW

	}

	//Attempt to load up the single user info
	p.User = &models.User{}
	err = user.ReadUser(p.User)
	if err != nil {
		fmt.Printf("Failed to load user info\n")
		p.ActiveView = REGISTER_VIEW
	}

	//Provide the scaffold view
	err = p.Controller("UtilityController").Endpoint("Scaffold", nil, session, p)

	return err
}

func (p *Application) MapViewKeys(key int) string {
	switch key {
	case MAIN_VIEW:
		return "MAIN_VIEW"
	case PROVISION_VIEW:
		return "PROVISION_VIEW"
	case DEVICE_DETAILS_VIEW:
		return "DEVICE_DETAILS_VIEW"
	case REGISTER_VIEW:
		return "REGISTER_VIEW"
	case LOG_VIEW:
		return "LOG_VIEW"
	case SUCCESS_VIEW:
		return "SUCCESS_VIEW"
	case LOAD_VIEW:
		return "LOAD_VIEW"
	case FATAL_VIEW:
		return "FATAL_VIEW"
	}
	return "UNKNOWN"
}

func (p *Application) NewSession() *Session {
	//generate unique session id
	session := NewSession("Provision", p.GetKeys())
	return session
}

func (p *Application) Close() {
	p.Connection.Close()
	p.BiChan <- EXIT
}

func (p *Application) GetChannel() chan int {
	return p.BiChan
}

func (p *Application) SetWindow(window *astilectron.Window) {
	p.Window = window
}

// The application should maintain a list of connected devices and pass on to the sessions as needed
func (p *Application) DeviceConnect(device *models.Device) error {
	var err error
	if device == nil {
		return fmt.Errorf("Device has nil handle, application can't connect\n")
	}
	ip := net.ByteToIP(device.IP)
	port := common.Int16(device.Port, common.LITTLE_ENDIAN)
	p.Connection, err = net.NewTCPConnection(ip.String() + ":" + fmt.Sprintf("%d", port))
	return err
}

func (p *Application) GetActiveView() string {
	return p.MapViewKeys(p.ActiveView)
}

func (p *Application) CloseDevice() {
	p.Connection.Close()
}

func (p *Application) GetKeys() map[string]*Session {
	return p.Sessions
}

func (p *Application) AddController(key string, controller Controller) error {
	p.Controllers[key] = controller
	return nil
}

func (p *Application) AddTemplate(key string, tmpl *templates.ApplicationTemplate) error {
	if tmpl.Template == nil {
		fmt.Printf("Failed to add template %s\n", key)
		return fmt.Errorf("Failed to add template %s\n", key)
	}
	p.TemplateViews[key] = tmpl
	return nil
}

func (p *Application) GetNet() *net.TCPConnection {
	return p.Connection
}

func (p *Application) GetElectron() *astilectron.Window {
	return p.Window
}

func (p *Application) GetSession(key string) *Session {
	return p.Sessions[key]
}

func (p *Application) GetTemplate(key string) *templates.ApplicationTemplate {
	return p.TemplateViews[key]
}

func (p *Application) Controller(key string) Controller {
	controller := p.Controllers[key]
	if controller == nil {
		fmt.Printf("Failed to find controller %s\n", key)
		return NewControllerBase()
	}
	return controller
}

func (p *Application) GenUniqueLocation() (net.IP, int) {
	netIP, port := GenIPLocation()
	for _, device := range p.Devices {
		if string(device.IP) == string(netIP.To4()) && int(common.Int32(device.Port, common.LITTLE_ENDIAN)) == port {
			netIP, port = GenIPLocation()
		}
	}
	return netIP, port
}

func GenIPLocation() (net.IP, int) {
	//IP Address Class C:
	var ip uint32
	var port uint16
	ip = 3232235520 + uint32(rand.Int31n(65535))
	port = 8000
	netIP := net.Int32ToIP(ip)
	return net.IP{IP: netIP}, int(port)
}
