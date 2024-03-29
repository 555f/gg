// Code generated by GG version dev. DO NOT EDIT.

//go:build !gg
// +build !gg

package server

import (
	"context"
	controller "github.com/555f/gg/examples/grpc-service/internal/usecase/controller"
	dto "github.com/555f/gg/examples/grpc-service/pkg/dto"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type options struct {
	profileController controller.ProfileController
}
type Option func(*options)

func ProfileController(s controller.ProfileController) Option {
	return func(o *options) {
		o.profileController = s
	}
}

type routeProfileController struct {
	UnimplementedProfileControllerServer
	svc controller.ProfileController
}

func (r *routeProfileController) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	profile, err := r.svc.Create(req.Token, req.FirstName, req.LastName, req.Address, int(req.Old), req.Age.AsTime(), req.Sleep.AsDuration())
	if err != nil {
		return nil, err
	}
	var resp *CreateResponse
	resp.Profile = &Profile{Id: int64(profile.ID), FistName: profile.FistName, LastName: profile.LastName, Address: &Address{Street: profile.Address.Street, Apt: int64(profile.Address.Apt), Apt2: int32(profile.Address.Apt2), Apt3: uint32(profile.Address.Apt3)}}
	return resp, err
}
func (r *routeProfileController) Remove(ctx context.Context, req *RemoveRequest) (*emptypb.Empty, error) {
	err := r.svc.Remove(req.Id)
	if err != nil {
		return nil, err
	}
	return nil, err
}
func (r *routeProfileController) Stream(stream ProfileController_StreamServer) error {
	chIn := make(chan *dto.Profile)
	go func() {
		for {
			data, err := stream.Recv()
			if err != nil {
				return
			}
			chIn <- &dto.Profile{ID: int(data.Id), FistName: data.FistName, LastName: data.LastName, Address: dto.Address{Street: data.Address.Street, Apt: int(data.Address.Apt), Apt2: int8(data.Address.Apt2), Apt3: uint8(data.Address.Apt3)}}
		}
	}()
	statistics, err := r.svc.Stream(chIn)
	for data := range statistics {
		stream.Send(&Statistic{ProfileID: int64(data.ProfileID), Sum: data.Sum, CreatedAt: timestamppb.New(data.CreatedAt)})
	}
	if err != nil {
		return err
	}
	return err
}
func (r *routeProfileController) Stream2(stream ProfileController_Stream2Server) error {
	chIn := make(chan *dto.Profile)
	go func() {
		for {
			data, err := stream.Recv()
			if err != nil {
				return
			}
			chIn <- &dto.Profile{ID: int(data.Id), FistName: data.FistName, LastName: data.LastName, Address: dto.Address{Street: data.Address.Street, Apt: int(data.Address.Apt), Apt2: int8(data.Address.Apt2), Apt3: uint8(data.Address.Apt3)}}
		}
	}()
	err := r.svc.Stream2(chIn)
	if err != nil {
		return err
	}
	return err
}
func (r *routeProfileController) Stream3(req *Stream3Request, stream ProfileController_Stream3Server) error {
	statistics, err := r.svc.Stream3(&dto.Profile{ID: int(req.Profile.Id), FistName: req.Profile.FistName, LastName: req.Profile.LastName, Address: dto.Address{Street: req.Profile.Address.Street, Apt: int(req.Profile.Address.Apt), Apt2: int8(req.Profile.Address.Apt2), Apt3: uint8(req.Profile.Address.Apt3)}})
	for data := range statistics {
		stream.Send(&Statistic{ProfileID: int64(data.ProfileID), Sum: data.Sum, CreatedAt: timestamppb.New(data.CreatedAt)})
	}
	if err != nil {
		return err
	}
	return err
}
func (r *routeProfileController) Update(ctx context.Context, req *UpdateRequest) (*emptypb.Empty, error) {
	err := r.svc.Update(dto.Profile{ID: int(req.Profile.Id), FistName: req.Profile.FistName, LastName: req.Profile.LastName, Address: dto.Address{Street: req.Profile.Address.Street, Apt: int(req.Profile.Address.Apt), Apt2: int8(req.Profile.Address.Apt2), Apt3: uint8(req.Profile.Address.Apt3)}})
	if err != nil {
		return nil, err
	}
	return nil, err
}
func Register(srv *grpc.Server, opts ...Option) {
	o := &options{}
	for _, f := range opts {
		f(o)
	}
	if o.profileController != nil {
		RegisterProfileControllerServer(srv, &routeProfileController{svc: o.profileController})
	}
}
