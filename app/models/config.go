package models

import "github.com/andewx/microxt/common"

type AppConfig struct {
	//Name of the application and some meta data
	AppName         string
	AppVersion      string
	AppEmail        string
	AppIcon         string
	AppLogo         string
	AppColor        string
	AppTheme        string
	AppDebug        bool
	AppPort         int
	ProvisionIP     string
	ProvisionPort   int
	DeviceIP        string
	DevicePort      int
	ApplicationPath string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{AppName: "AirXT", AppVersion: "0.0.1", AppEmail: "feedback@dieselx.io", AppIcon: "microxt.png", AppLogo: "microxt.png", AppColor: "#000000", AppTheme: "dark", AppDebug: true, ProvisionIP: "192.168.0.10",
		AppPort: 9060, DeviceIP: "", DevicePort: 9060, ApplicationPath: common.ProjectRelativePath("microxt/app")}
}

type RadarData struct {
	ADCData  *ADCData  `json:"adcdata"`
	FFTData  *FFTData  `json:"fftdata"`
	PDATData *PDATData `json:"pdatdata"`
	TDATData *TDATData `json:"tdatdata"`
	DDATData *DDATData `json:"ddatdata"`
}

func NewRadarData() *RadarData {
	return &RadarData{ADCData: NewADCData(), FFTData: NewFFTData(), PDATData: NewPDATData(), TDATData: NewTDATData(), DDATData: NewDDATData()}
}

type SessionObject struct {
	FirstName     string      `json:"firstname"`
	LastName      string      `json:"lastname"`
	Email         string      `json:"email"`
	UserID        int64       `json:"userid"`
	ActiveView    string      `json:"activeview"`
	ActiveToolbar string      `json:"activetoolbar"`
	ActiveDevice  string      `json:"activedevice"`
	ActiveNav     string      `json:"activenav"`
	Terminal      string      `json:"terminal"`
	ParamsRadar   RadarParams `json:"paramsradar"`
	RadarData     *RadarData  `json:"radardata"`
	Devices       []*Device   `json:"devices"`
	//private
	logging *Logging
}

func NewSessionObject() *SessionObject {
	return &SessionObject{ActiveView: "Provision", ActiveToolbar: "config", ActiveDevice: "none", FirstName: "DieselX", LastName: "User", ActiveNav: "config", Terminal: "shell", ParamsRadar: RadarParams{}, RadarData: NewRadarData(), Devices: make([]*Device, 0), logging: NewLogging()}
}

type Device struct {
	ID     int64  `json:"deviceid"`
	Name   string `json:"devicename"`
	IP     string `json:"deviceip"`
	Port   int    `json:"deviceport"`
	Driver string `json:"devicedriver"`
	Status int    `json:"devicestatus"`
}
