package templates

import (
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
	templ.Create(name, filenames...)
	return templ
}

func (t *ApplicationTemplate) Create(name string, filenames ...string) error {
	var err error
	t.Template = template.New(name)
	t.Name = name
	t.Template, err = template.ParseFiles(filenames...)
	return err
}

func (t *ApplicationTemplate) Execute(obj *models.SessionObject, w io.Writer) error {
	var err error
	ts := t.Template.Templates()
	for _, tmp := range ts {
		if err = tmp.Execute(w, obj); err != nil {
			return err
		}
	}

	return err
}
