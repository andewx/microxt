package templates

import (
	"text/template"
)

type Template interface {
	Info() string
	AddItem(interface{}, string) error
	Get() string
}

type ProvisionTemplate struct {
	Template *template.Template
}

func (t *ProvisionTemplate) Create() error {
	var err error
	t.Template, err = template.New("Provision").ParseFiles("provision.html")
	return err
}

type IDETemplate struct {
	Template *template.Template
}

func (t *IDETemplate) Create() error {
	var err error
	t.Template, err = template.New("IDE").ParseFiles("ide.html")
	return err
}

type LoginTemplate struct {
	Template *template.Template
}

func (t *LoginTemplate) Create() error {
	var err error
	t.Template, err = template.New("IDE").ParseFiles("ide.html")
	return err
}
