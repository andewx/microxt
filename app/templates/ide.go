package templates

import "github.com/andewx/microxt/app/models"

type Tab struct {
	Name   string
	Active bool
	Link   string
}

type Device struct {
	Name string
	ID   string
	IP   string
	Port string
}

type Timeline struct {
	X       int
	Y       int
	Width   int
	Height  int
	Active  bool
	Scroll  float32
	MaxRSSI []int
}

type Terminal struct {
	Current string
	Buffer  []string
}

type IDE struct {
	Devices      []*Device
	ActiveDevice *Device
	NavTabs      []*Tab
	FooterTabs   []*Tab
	Session      *models.SessionObject
}
