// @gg:"test"

package controller

import (
	"context"
	"time"

	"github.com/555f/gg/examples/grpc-service/pkg/dto"
)

type ContextKey int

const (
	TestMetaContextKey ContextKey = iota + 1
)

// ProfileController Профиль пользователя
// Методы для работы с профилем пользователя
// @gg:"grpc"
// @gg:"middleware"
// @gg:"klog"
// @grpc-server
// @grpc-client
type ProfileController interface {
	// Create Создать профиль
	// @grpc-meta-context:"~/examples/grpc-service/internal/usecase/controller.TestMetaContextKey"
	Create(
		ctx context.Context,
		// @grpc-version:"1"
		token string,
		// @grpc-version:"2"
		firstName string,
		// @grpc-version:"3"
		lastName string,
		// @grpc-version:"4"
		address string,
		// @grpc-version:"5"
		old int,
		// @grpc-version:"6"
		age time.Time,
		// @grpc-version:"7"
		sleep time.Duration,
	) (
		// @grpc-version:"1"
		profile *dto.Profile,
		err error,
	)
	Update(
		// @grpc-version:"1"
		profile dto.Profile,
	) (err error)
	// Remove Удалить профиль
	Remove(
		// @grpc-version:"1"
		id string,
	) (
		err error,
	)
	Stream(profile chan *dto.Profile) (statistics chan *dto.Statistic, err error)
	Stream2(profile chan *dto.Profile) (err error)
	Stream3(
		// @grpc-version:"1"
		profile *dto.Profile,
	) (statistics chan *dto.Statistic, err error)
}
