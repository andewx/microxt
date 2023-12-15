package application

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andewx/microxt/app/models"
	"github.com/andewx/microxt/common"
)

type UserController struct {
	*ControllerBase
}

func NewUserController() *UserController {
	c := &UserController{NewControllerBase()}
	c.handles["Register"] = c.Register
	return c
}

func (u *UserController) ReadUser(user *models.User) error {
	//Reads devices.json file and returns the contents as a string
	filename := common.ProjectRelativePath("microxt/user.json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, user)
	if err != nil {
		return err
	}

	fmt.Printf("Read user.json file successfully\n")
	return nil
}

func (u *UserController) SaveUser(user *models.User) error {
	bytes, err := json.Marshal(user)
	filename := common.ProjectRelativePath("microxt/user.json")
	if err == nil {
		file, err_ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		defer file.Close()
		if err_ == nil {
			file.Write(bytes)
		} else {
			fmt.Printf("Failed to open user.json file %s\n", err_.Error())
		}
	} else {
		fmt.Printf("Failed to marshal user %s\n", err.Error())
	}
	return err
}

func (u *UserController) Register(params map[string]string, session *Session, app *Application) error {
	first := params["first"]
	last := params["last"]
	email := params["email"]
	password := params["password"]
	app.User = &models.User{First: first, Last: last, Email: email, Password: password}
	if err := u.SaveUser(app.User); err != nil {
		app.Controller("UtilityController").Endpoint("SendNotification", map[string]string{"message": "Failed to save user"}, session, app)
		return err
	} else {
		fmt.Printf("Saved user %s\n", app.User.Email)
		app.ActiveView = PROVISION_VIEW
		app.Controller("UtilityController").Endpoint("Scaffold", nil, session, app)
	}
	return nil
}
