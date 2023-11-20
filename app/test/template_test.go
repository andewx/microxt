package test

import (
	"os"
	"testing"
	"text/template"

	"github.com/andewx/microxt/app/models"
	"github.com/andewx/microxt/app/templates"
	"github.com/andewx/microxt/common"
)

func TestTemplate(t *testing.T) {
	var err error
	tmpl := new(templates.ApplicationTemplate)
	file1 := common.ProjectRelativePath("app/templates/test.gohtml")
	file2 := common.ProjectRelativePath("app/templates/login.gohtml")
	var session *models.SessionObject
	session = models.NewSessionObject()
	err = tmpl.Create("test", file1, file2)

	if err != nil {
		t.Errorf("Failed to create template")
	}

	err = tmpl.Execute(session, os.Stdout)

	if err != nil {
		t.Errorf("Failed to execute template %s", err)
	}

	tmplT, err := template.New("test").Parse("<div><p>Hello, {{.FirstName}} {{.LastName}}</p></div>")
	if err != nil {
		t.Errorf("Base template failed\n")
		return
	}

	err = tmplT.Execute(os.Stdout, session)

	if err != nil {
		t.Errorf("Failed to execute template %s", err)
		return
	}

}
