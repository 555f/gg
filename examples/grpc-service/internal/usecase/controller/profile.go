// @gg:"test"

package controller

import (
	"github.com/555f/gg/examples/grpc-service/pkg/dto"
)

// ProfileController Профиль пользователя
// Методы для работы с профилем пользователя
// @gg:"grpc"
// @grpc-server
// @grpc-interface:"~/examples/grpc-service/internal/interface/controller/pb.GreeterServer"
type ProfileController interface {
	// Create Создать профиль
	// @jsonrpc-name:"profile.create"
	Create(
		token string,
		firstName string,
		lastName string,
		address string,
	) (
		profile *dto.Profile,
		err error,
	)
	// Remove Удалить профиль
	// @jsonrpc-name:"profile.delete"
	Remove(
		id string,
	) (
		err error,
	)
}
