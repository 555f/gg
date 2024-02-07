// Code generated by GG version . DO NOT EDIT.

//go:build !gg
// +build !gg

package middleware

import (
	controller "github.com/555f/gg/examples/rest-service-echo/internal/usecase/controller"
	dto "github.com/555f/gg/examples/rest-service-echo/pkg/dto"
)

type ProfileControllerMiddleware func(controller.ProfileController) controller.ProfileController

func ProfileControllerMiddlewareChain(outer ProfileControllerMiddleware, others ...ProfileControllerMiddleware) ProfileControllerMiddleware {
	return func(next controller.ProfileController) controller.ProfileController {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}

type profileControllerBaseMiddleware struct {
	next     controller.ProfileController
	mediator any
}

func (m *profileControllerBaseMiddleware) Create(firstName string, lastName string, address string, zip int) (profile *dto.Profile, err error) {
	defer func() {
		if s, ok := m.mediator.(profileControllerCreateBaseMiddleware); ok {
			s.Create(firstName, lastName, address, zip)
		}
	}()
	return m.next.Create(firstName, lastName, address, zip)
}
func (m *profileControllerBaseMiddleware) DownloadFile(id string) (data string, err error) {
	defer func() {
		if s, ok := m.mediator.(profileControllerDownloadFileBaseMiddleware); ok {
			s.DownloadFile(id)
		}
	}()
	return m.next.DownloadFile(id)
}
func (m *profileControllerBaseMiddleware) Remove(id string) (err error) {
	defer func() {
		if s, ok := m.mediator.(profileControllerRemoveBaseMiddleware); ok {
			s.Remove(id)
		}
	}()
	return m.next.Remove(id)
}

type profileControllerCreateBaseMiddleware interface {
	Create(firstName string, lastName string, address string, zip int)
}
type profileControllerDownloadFileBaseMiddleware interface {
	DownloadFile(id string)
}
type profileControllerRemoveBaseMiddleware interface {
	Remove(id string)
}

func ProfileControllerBaseMiddleware(mediator any) ProfileControllerMiddleware {
	return func(next controller.ProfileController) controller.ProfileController {
		return &profileControllerBaseMiddleware{next: next, mediator: mediator}
	}
}