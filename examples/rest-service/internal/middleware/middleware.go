package middleware

import controller "github.com/f555/gg-examples/internal/usecase/controller"

type ProfileControllerMiddleware func(controller.ProfileController) controller.ProfileController

func ProfileControllerMiddlewareChain(outer ProfileControllerMiddleware, others ...ProfileControllerMiddleware) ProfileControllerMiddleware {
	return func(next controller.ProfileController) controller.ProfileController {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}
