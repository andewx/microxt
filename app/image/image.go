package image

import "github.com/andewx/microxt/app/models"

type Image interface {
	GetItemByID(uuid int) (models.Model, error)
	GetItemByKey(key string) (models.Model, error)
	GetItemByGroup(group string) ([]models.Model, error)
	AddItem(models.Model) error
	SaveImage(url string) error
	LoadImage(url string) error
}

type image struct {
	objects map[string]models.Model
}
