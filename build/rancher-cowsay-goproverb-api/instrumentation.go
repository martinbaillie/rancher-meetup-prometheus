package main

import (
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
)

func instrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	sayResult metrics.Gauge,
) ServiceMiddleware {
	return func(next Service) Service {
		return instrmw{requestCount, requestLatency, sayResult, next}
	}
}

type instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	sayResult      metrics.Gauge
	Service
}

func (mw instrmw) Textsay(reqindex int) (retindex int, say string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "textsay", "index", strconv.Itoa(retindex)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.sayResult.With("method", "textsay").Set(float64(retindex))
	}(time.Now())

	retindex, say = mw.Service.Textsay(reqindex)
	return
}

func (mw instrmw) Cowsay(reqindex int) (retindex int, say string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "cowsay", "index", strconv.Itoa(retindex)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.sayResult.With("method", "cowsay").Set(float64(retindex))
	}(time.Now())

	retindex, say = mw.Service.Cowsay(reqindex)
	return
}
