package controller

import (
	"github.com/555f/gg/examples/rest-service/internal/usecase/controller"
	"github.com/555f/gg/examples/rest-service/pkg/dto"
)

var _ controller.ProfileController = &ProfileController{}

type ProfileController struct{}

func (c *ProfileController) DownloadFile(id string) (data string, err error) {
	return "data", nil
}

func (c *ProfileController) Create(firstName string, lastName string, address string, zip int) (profile *dto.Profile, err error) {
	return &dto.Profile{
		FistName: "Vitaly",
		LastName: "Lobchuk",
		Address:  dto.Address{},
	}, nil
}

func (c *ProfileController) Remove(id string) (err error) {
	return nil
}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}
