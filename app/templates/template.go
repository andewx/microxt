package templates

import (
	"fmt"
	"html/template"
	"io"

	"github.com/andewx/microxt/app/models"
)

type StringWriter struct {
	Str string
}

func (s *StringWriter) Write(p []byte) (int, error) {
	s.Str = string(p)
	return len(p), nil
}

//--------------Template---------------//

type Template interface {
	Create(string, ...string) error
	Execute(obj *models.SessionObject, w io.Writer) error
}

type ApplicationTemplate struct {
	Template *template.Template
	Name     string
}

func NewTemplate(name string, filenames ...string) *ApplicationTemplate {
	templ := &ApplicationTemplate{}
	err := templ.Create(name, filenames...)
	if err != nil {
		fmt.Printf("Failed to create template %s with filename%s\n", name, filenames)
	}
	return templ
}

func (t *ApplicationTemplate) Create(name string, filenames ...string) error {
	var err error
	t.Template = template.New(name)
	t.Name = name
	t.Template, err = template.ParseFiles(filenames...)
	return err
}

func (t *ApplicationTemplate) Execute(obj any, w io.Writer) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Execute", r)
			err = fmt.Errorf("error: %v", r)
		}
	}()

	//Recover from panic
	ts := t.Template.Templates()
	for _, tmp := range ts {
		tmp.Execute(w, obj)
	}
	return err
}
