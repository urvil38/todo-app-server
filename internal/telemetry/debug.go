package telemetry

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"

	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/go-chi/chi/v5"
	prom_client "github.com/prometheus/client_golang/prometheus"
	"github.com/urvil38/todo-app/internal/config"
	"github.com/urvil38/todo-app/internal/version"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/runmetrics"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
)

// RouteTagger is a func that can be used to derive a dynamic route tag for an
// incoming request.
type RouteTagger func(route string, r *http.Request) string

// Router is an http multiplexer that instruments per-handler debugging
// information and census instrumentation.
type Router struct {
	http.Handler
	mux    *chi.Mux
	tagger RouteTagger
}

// NewRouter creates a new Router, using tagger to tag incoming requests in
// monitoring. If tagger is nil, a default route tagger is used.
func NewRouter(tagger RouteTagger) *Router {
	if tagger == nil {
		tagger = func(route string, r *http.Request) string {
			return strings.Trim(route, "/")
		}
	}
	router := chi.NewRouter()
	return &Router{
		mux:     router,
		Handler: &ochttp.Handler{Handler: router},
		tagger:  tagger,
	}
}

// Handle registers handler with the given route. It has the same routing
// semantics as http.ServeMux.
func (r *Router) Handle(method, route string, handler http.Handler) {
	r.mux.MethodFunc(method, route, func(w http.ResponseWriter, req *http.Request) {
		tag := r.tagger(route, req)
		ochttp.WithRouteTag(handler, tag).ServeHTTP(w, req)
	})
}

const debugPage = `
<html>
<p><a href="/tracez">/tracez</a> - trace spans</p>
<p><a href="/statsz">/statz</a> - prometheus metrics page</p>
`

// Init configures tracing and aggregation according to the given Views.
func Init(cfg config.Config, views ...*view.View) error {
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	if err := view.Register(views...); err != nil {
		return fmt.Errorf("debug.Init(views): view.Register: %v", err)
	}

	exp, err := ocagent.NewExporter(ocagent.WithInsecure(), ocagent.WithServiceName("todo-app"))
	if err != nil {
		log.Fatalf("Failed to create the agent exporter: %v", err)
	}

	trace.RegisterExporter(exp)

	err = runmetrics.Enable(runmetrics.RunMetricOptions{
		EnableCPU:            true,
		EnableMemory:         true,
		UseDerivedCumulative: true,
		Prefix:               "todo_app/",
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// NewServer creates a new http.Handler for serving debug information.
func NewServer(cfg config.Config) (http.Handler, error) {
	pe, err := prometheus.NewExporter(prometheus.Options{
		ConstLabels: prom_client.Labels{
			"service":     "todo-server",
			"environment": cfg.Env,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("debug.NewServer: prometheus.NewExporter: %v", err)
	}
	mux := http.NewServeMux()
	zpages.Handle(mux, "/")
	mux.Handle("/statsz", pe)
	mux.Handle("/version", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf("version: %v\ncommit: %v", version.Version, version.Commit))
	}))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, debugPage)
	})
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	return mux, nil
}
