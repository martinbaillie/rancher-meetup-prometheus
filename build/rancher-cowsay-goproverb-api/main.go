package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	stdlog "log"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// Linker-provided project/build information
var (
	projectName,
	projectVersion,
	buildTime,
	buildHash,
	buildUser string
)

func main() {
	// Behaviour
	const (
		defHTTPBasePath   = "/v1"
		defHTTPPort       = 8080
		defPrometheusPort = 8081
	)
	var (
		// In keeping with 12 factor, flags can also be set in the environment.
		// Do this by uppercasing the entire CLI flag e.g. HTTP_PORT.
		debug          = flag.Bool("debug", false, "Turn on debug logging output")
		addr           = flag.String("addr", defaultAddr(), "Service transport address")
		httpBasepath   = flag.String("http_basepath", defHTTPBasePath, "Basepath to serve the HTTP endpoints from")
		httpPort       = flag.Int("http_port", defHTTPPort, "HTTP transport bind port")
		prometheusPort = flag.Int("prometheus_port", defPrometheusPort, "Metrics (Prometheus) bind port")
	)
	flag.Parse()

	var (
		httpAddr       = fmt.Sprintf("%s:%d", *addr, *httpPort)
		prometheusAddr = fmt.Sprintf("%s:%d", *addr, *prometheusPort)
	)

	// Logging
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)

		if *debug {
			// Show debug level log statements
			logger = level.NewFilter(logger, level.AllowDebug())
		} else {
			logger = level.NewFilter(logger, level.AllowInfo())
		}

		// Redirect stdlib logger to Go kit logger.
		stdlog.SetOutput(log.NewStdlibAdapter(logger))

		logger = log.With(logger, "caller", log.DefaultCaller)
		logger = log.With(logger, "ts", log.DefaultTimestamp)

		logBuildInfo(logger, debug)
	}
	defer logger.Log("service", projectName, "msg", "stopping")

	// Context plumbing and interrupt/error channels
	//ctx := context.Background()
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Instrumentation (Prometheus)
	requestCount := prometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name: "rancher_cowsay_goproverb_api_request_count",
		Help: "Number of requests received.",
	}, []string{"method", "index"})
	requestLatency := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Name: "rancher_cowsay_goproverb_api_request_latency_seconds",
		Help: "Total duration of requests in seconds.",
	}, []string{"method", "index"})
	sayResult := prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Name: "rancher_cowsay_goproverb_say_result_index",
		Help: "The current index result of each (text|cow)say method.",
	}, []string{"method"})

	var svc Service
	{
		svc = service{}
		svc = loggingMiddleware(logger)(svc)
		svc = instrumentingMiddleware(requestCount, requestLatency, sayResult)(svc)
	}

	// HTTP transport
	go func() {
		logger := log.With(logger, "service", projectName)

		// Create the router
		r := mux.NewRouter().StrictSlash(true)

		// Log every 404, even those not covered by our pre-defined routes
		r.NotFoundHandler = notFoundLogger(logger)

		// Further decorate the router with useful HTTP middlewares
		var rmws http.Handler = r
		rmws = handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(rmws)
		rmws = handlers.CompressHandler(rmws)
		rmws = handlers.ProxyHeaders(rmws)
		rmws = handlers.RecoveryHandler(handlers.RecoveryLogger(wrapLogger{logger}))(rmws)

		// Make our service HTTP handlers
		r.Handle(defHTTPBasePath+"/textsay", httptransport.NewServer(
			makeTextsayEndpoint(svc),
			decodeRequest,
			encodeGenericJSONResponse,
		))
		r.Handle(defHTTPBasePath+"/cowsay", httptransport.NewServer(
			makeCowsayEndpoint(svc),
			decodeRequest,
			encodeCowsayTextResponse,
		))

		logger.Log("msg", "listening", "addr", httpAddr, "base_path", *httpBasepath)
		errc <- http.ListenAndServe(httpAddr, rmws)
	}()

	// Metrics (Prometheus) transport
	go func() {
		logger := log.With(logger, "metrics", "Prometheus")

		r := mux.NewRouter()
		r.Handle("/metrics", promhttp.Handler())

		logger.Log("msg", "listening", "addr", prometheusAddr, "base_path", "/metrics")
		errc <- http.ListenAndServe(prometheusAddr, r)
	}()

	// Run!
	logger.Log("service", projectName, "msg", <-errc)
}

func defaultAddr() string {
	// Try to get the IP used in the default route
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "0.0.0.0"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}

// Useful error logging helpers
func notFoundLogger(logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log("err", http.StatusText(http.StatusNotFound), "url", r.URL)
		w.WriteHeader(http.StatusNotFound)
	})
}

// wrapLogger wraps a Go kit logger so we can use it as the logging service for
// Gorilla middlewares like the recovery handler.
type wrapLogger struct {
	log.Logger
}

func (logger wrapLogger) Println(args ...interface{}) {
	logger.Log("msg", fmt.Sprint(args...))
}

func logBuildInfo(logger log.Logger, debug *bool) {
	logger.Log("service", projectName, "msg", "starting",
		"version", projectVersion, "debug", *debug)

	logger.Log("build_time", buildTime, "build_commit", buildHash,
		"build_user", buildUser)

}
