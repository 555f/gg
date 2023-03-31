// @gg:"test"

package controller

import (
	"github.com/f555/gg-examples/pkg/dto"
)

// ProfileController Профиль пользователя
// Методы для работы с профилем пользователя
// @gg:"http"
// @gg:"middleware"
// @gg:"logging"
// @http-api-doc
// @http-type:"jsonrpc"
// @http-client:"pkg/client/client.go"
// @http-server:"internal/transport/server.go"
type ProfileController interface {
	// Create Создать профиль
	// @http-path:"profile.create"
	Create(
		// @http-required
		firstName string,
		// @http-required
		lastName string,
		address string,
	) (
		profile *dto.Profile,
		err error,
	)
	// Remove Удалить профиль
	// @http-path:"profile.delete"
	Remove(
		// @http-required
		id string,
	) (
		err error,
	)
}
