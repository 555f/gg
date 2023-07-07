// @gg:"test"

package controller

import (
	"github.com/555f/gg/examples/jsonrpc-service/pkg/dto"
)

// ProfileController Профиль пользователя
// Методы для работы с профилем пользователя
// @gg:"http"
// @gg:"middleware"
// @gg:"logging"
// @http-api-doc
// @http-type:"jsonrpc"
// @http-client
// @http-server
type ProfileController interface {
	// Create Создать профиль
	// @http-path:"profile.create"
	Create(
		// @http-required
		token string,
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
