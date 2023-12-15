package models

type User struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionObject struct {
	User    *User `json:"user"`
	logging *Logging
}

func NewSessionObject() *SessionObject {
	return &SessionObject{logging: NewLogging()}
}
