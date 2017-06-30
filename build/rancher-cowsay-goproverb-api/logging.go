package main

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	Service
}

func (mw logmw) Textsay(reqindex int) (retindex int, say string) {
	defer func(begin time.Time) {
		_ = level.Debug(mw.logger).Log(
			"method", "textsay",
			"index", retindex,
			"say", say,
			"took", time.Since(begin),
		)
	}(time.Now())

	retindex, say = mw.Service.Textsay(reqindex)
	return
}

func (mw logmw) Cowsay(reqindex int) (retindex int, say string) {
	defer func(begin time.Time) {
		_ = level.Debug(mw.logger).Log(
			"method", "cowsay",
			"index", retindex,
			"took", time.Since(begin),
		)
	}(time.Now())

	retindex, say = mw.Service.Cowsay(reqindex)
	return
}
