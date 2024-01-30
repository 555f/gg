package controller

import (
	"github.com/555f/gg/examples/rest-service/pkg/dto"
)

// ProfileController Профиль пользователя
// Методы для работы с профилем пользователя
// @gg:"http"
// @gg:"middleware"
// @gg:"klog"
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
	// @http-content-type:"xml,root=profile"
	// @http-content-type:"urlencoded"
	// @http-content-type:"multipart"
	// @http-accept-type:"json"
	// @http-accept-type:"xml"
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
	// @http-path:"/profiles/:id"
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
