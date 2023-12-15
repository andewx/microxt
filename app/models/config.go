package models

import (
	"github.com/andewx/microxt/common"
)

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
