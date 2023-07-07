package controller

import (
	"github.com/555f/gg/examples/jsonrpc-service/internal/usecase/controller"
	"github.com/555f/gg/examples/jsonrpc-service/pkg/dto"
)

var _ controller.ProfileController = &ProfileController{}

type ProfileController struct{}

func (p *ProfileController) Create(token string, firstName string, lastName string, address string) (profile *dto.Profile, err error) {
	return &dto.Profile{
		FistName: "Vitaly",
		LastName: "Lobchuk",
		Address:  dto.Address{},
	}, nil
}

func (p *ProfileController) Remove(id string) (err error) {
	return nil
}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}
