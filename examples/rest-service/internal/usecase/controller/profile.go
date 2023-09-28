package controller

import (
	"github.com/555f/gg/examples/rest-service/pkg/dto"
)

// ProfileController Профиль пользователя
// Методы для работы с профилем пользователя
// @gg:"http"
// @gg:"middleware"
// @gg:"logging"
// @http-type:"echo"
// @http-openapi
// @http-api-doc
// @http-openapi-tags:"profile"
// @http-client
// @http-server
type ProfileController interface {
	// Create Создать профиль
	// @http-path:"/profiles"
	// @http-method:"POST"
	// @http-content-types:"json,xml,urlencoded,multipart,root-xml=profile"
	// @http-accept-types:"json"
	Create(
		// @http-required
		firstName string,
		// @http-required
		lastName string,
		address string,
		// @http-required
		zip int,
	) (
		profile *dto.Profile,
		err error,
	)
	// Remove Удалить профиль
	// @http-path:"/profiles/{id}"
	// @http-method:"DELETE"
	Remove(
		// @http-required
		id string,
	) (
		err error,
	)
	// DownloadFile
	// @http-path:"/profiles/:id/file"
	// @http-method:"GET"
	DownloadFile(
		// @http-required
		id string,
	) (
		data string,
		err error,
	)
}
