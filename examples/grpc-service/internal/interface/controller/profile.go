package controller

import (
	"fmt"
	"time"

	"github.com/555f/gg/examples/grpc-service/internal/usecase/controller"
	"github.com/555f/gg/examples/grpc-service/pkg/dto"
)

var _ controller.ProfileController = &ProfileController{}

type ProfileController struct{}

// Stream3 implements controller.ProfileController.
func (*ProfileController) Stream3(profile *dto.Profile) (statistics chan *dto.Statistic, err error) {
	panic("unimplemented")
}

// Stream2 implements controller.ProfileController.
func (*ProfileController) Stream2(profile chan *dto.Profile) (err error) {
	go func() {
		for p := range profile {
			fmt.Printf("%v\n", p)
		}
	}()
	return nil
}

// Stream implements controller.ProfileController.
func (*ProfileController) Stream(profile chan *dto.Profile) (statistics chan *dto.Statistic, err error) {
	defer func() {
		fmt.Println("disconnect stream")
	}()
	outCh := make(chan *dto.Statistic)
	go func() {
		for p := range profile {
			outCh <- &dto.Statistic{
				ProfileID: p.ID,
				Sum:       100,
			}
		}
	}()
	return outCh, nil
}

// Create implements controller.ProfileController.
func (*ProfileController) Create(token string, firstName string, lastName string, address string, old int, age time.Time, sleep time.Duration) (profile *dto.Profile, err error) {
	return &dto.Profile{
		FistName: "Vitaly",
		LastName: "Lobchuk",
		Address:  dto.Address{},
	}, nil
}

// Update implements controller.ProfileController.
func (*ProfileController) Update(profile dto.Profile) (err error) {
	return nil
}

func (p *ProfileController) Remove(id string) (err error) {
	return nil
}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}
