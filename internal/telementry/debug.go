package telementry

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/prometheus"
	prom_client "github.com/prometheus/client_golang/prometheus"
	"github.com/gorilla/mux"
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
	mux    *mux.Router
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
	mux := mux.NewRouter()
	return &Router{
		mux:     mux,
		Handler: &ochttp.Handler{Handler: mux},
		tagger:  tagger,
	}
}

// Handle registers handler with the given route. It has the same routing
// semantics as http.ServeMux.
func (r *Router) Handle(route string, handler http.Handler) *mux.Route {
	return r.mux.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
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

	if cfg.TracingAgentURI != "" && cfg.TracingCollectorURI != "" {

		je, err := jaeger.NewExporter(jaeger.Options{
			AgentEndpoint:     cfg.TracingAgentURI,
			CollectorEndpoint: cfg.TracingCollectorURI,
			ServiceName:       "todo-app",
		})

		if err != nil {
			return fmt.Errorf("Failed to create the Jaeger exporter: %v", err)
		}

		trace.RegisterExporter(je)
	}

	exp, err := ocagent.NewExporter(ocagent.WithInsecure(), ocagent.WithServiceName("todo-app"))
	if err != nil {
		log.Fatalf("Failed to create the agent exporter: %v", err)
	}

	trace.RegisterExporter(exp)

	err = runmetrics.Enable(runmetrics.RunMetricOptions{
		EnableCPU:    true,
		EnableMemory: true,
		UseDerivedCumulative: true,
		Prefix:       "todo_app/",
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// NewServer creates a new http.Handler for serving debug information.
func NewServer() (http.Handler, error) {
	pe, err := prometheus.NewExporter(prometheus.Options{
		ConstLabels: prom_client.Labels{
			"version": version.Version,
			"rev": version.Commit,
			"service": "todo-server",
			"environment": "dev",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("debug.NewServer: prometheus.NewExporter: %v", err)
	}
	mux := http.NewServeMux()
	zpages.Handle(mux, "/")
	mux.Handle("/statsz", pe)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, debugPage)
	})

	return mux, nil
}
