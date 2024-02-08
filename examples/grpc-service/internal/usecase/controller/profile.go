// @gg:"test"

package controller

import (
	"github.com/555f/gg/examples/grpc-service/pkg/dto"
)

// ProfileController Профиль пользователя
// Методы для работы с профилем пользователя
// @gg:"grpc"
// @grpc-server
type ProfileController interface {
	// Create Создать профиль
	Create(
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
	) (
		// @grpc-version:"1"
		profile *dto.Profile,
		err error,
	)
	// Remove Удалить профиль
	Remove(
		// @grpc-version:"1"
		id string,
	) (
		err error,
	)
}
