package logging

import (
	middleware "github.com/555f/gg/examples/rest-service/internal/middleware"
	controller "github.com/555f/gg/examples/rest-service/internal/usecase/controller"
	dto "github.com/555f/gg/examples/rest-service/pkg/dto"
	log "github.com/go-kit/log"
	level "github.com/go-kit/log/level"
	"time"
)

type errLevel interface {
	Level() string
}
type logError interface {
	LogError() error
}

func levelLogger(e errLevel, logger log.Logger) log.Logger {
	switch e.Level() {
	default:
		return level.Error(logger)
	case "debug":
		return level.Debug(logger)
	case "info":
		return level.Info(logger)
	case "warn":
		return level.Warn(logger)
	}
}

type ProfileControllerLoggingMiddleware struct {
	next   controller.ProfileController
	logger log.Logger
}

func (s *ProfileControllerLoggingMiddleware) Create(firstName string, lastName string, address string, zip int) (profile *dto.Profile, err error) {
	defer func(now time.Time) {
		logger := log.WithPrefix(s.logger, "message", "call method - Create", "firstName", firstName, "lastName", lastName, "address", address, "zip", zip)
		if err != nil {
			if e, ok := err.(errLevel); ok {
				logger = levelLogger(e, logger)
			} else {
				logger = level.Error(logger)
			}
			if e, ok := err.(logError); ok {
				logger = log.WithPrefix(logger, "err", e.LogError())
			} else {
				logger = log.WithPrefix(logger, "err", err)
			}
		} else {
			logger = level.Debug(logger)
			logger = log.WithPrefix(logger, "profile", "")
		}
		_ = logger.Log("dur", time.Since(now))
	}(time.Now())
	profile, err = s.next.Create(firstName, lastName, address, zip)
	return
}
func (s *ProfileControllerLoggingMiddleware) DownloadFile(id string) (data string, err error) {
	defer func(now time.Time) {
		logger := log.WithPrefix(s.logger, "message", "call method - DownloadFile", "id", id)
		if err != nil {
			if e, ok := err.(errLevel); ok {
				logger = levelLogger(e, logger)
			} else {
				logger = level.Error(logger)
			}
			if e, ok := err.(logError); ok {
				logger = log.WithPrefix(logger, "err", e.LogError())
			} else {
				logger = log.WithPrefix(logger, "err", err)
			}
		} else {
			logger = level.Debug(logger)
			logger = log.WithPrefix(logger, "data", data)
		}
		_ = logger.Log("dur", time.Since(now))
	}(time.Now())
	data, err = s.next.DownloadFile(id)
	return
}
func (s *ProfileControllerLoggingMiddleware) Remove(id string) (err error) {
	defer func(now time.Time) {
		logger := log.WithPrefix(s.logger, "message", "call method - Remove", "id", id)
		if err != nil {
			if e, ok := err.(errLevel); ok {
				logger = levelLogger(e, logger)
			} else {
				logger = level.Error(logger)
			}
			if e, ok := err.(logError); ok {
				logger = log.WithPrefix(logger, "err", e.LogError())
			} else {
				logger = log.WithPrefix(logger, "err", err)
			}
		} else {
			logger = level.Debug(logger)
			logger = log.WithPrefix(logger)
		}
		_ = logger.Log("dur", time.Since(now))
	}(time.Now())
	err = s.next.Remove(id)
	return
}
func LoggingProfileControllerMiddleware(logger log.Logger) middleware.ProfileControllerMiddleware {
	return func(next controller.ProfileController) controller.ProfileController {
		return &ProfileControllerLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}
