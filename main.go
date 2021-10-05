package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/jtarte/sample_fluentd/utils"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//App application to run the server
type App struct {
	router *mux.Router
	reg    *prometheus.Registry
}

// main
// entry point of the program
func main() {
	//init the logger
	initLog()
	// define the URL handler
	a := App{}
	a.init()
	// launch the http server
	myport := os.Getenv("PORT")
	if myport == "" {
		myport = ":8080"
	} else {
		log.Trace("Found a custom port :" + myport)
		myport = ":" + myport
	}
	log.Info("Server started on " + myport)
	log.Fatal(http.ListenAndServe(myport, a.router))
}

// define the handler for the uri with instrumentation
// uri the uri to handle
// label the label of handler
// f the handler function
func (a *App) instrumentHandle(uri string, label string, f http.HandlerFunc) *mux.Route {
	log.Debug("defining the collected metric for " + uri)
	httpRequestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: label + "_http_requests_total",
		Help: "Count of all HTTP requests on /" + label,
	}, []string{"code", "method"})
	httpRequestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: label + "_http_request_duration_seconds",
		Help: "Duration of all HTTP requests of /" + label,
	}, []string{"code", "handler", "method"})
	a.reg.MustRegister(httpRequestsTotal)
	a.reg.MustRegister(httpRequestDuration)
	return a.router.Handle(uri, promhttp.InstrumentHandlerDuration(httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": label}), promhttp.InstrumentHandlerCounter(httpRequestsTotal, http.HandlerFunc(f))))
}

// Init init the http server
func (a *App) init() {
	log.Debug("Entering App.init")
	a.reg = prometheus.NewRegistry()
	//basic handling
	a.router = mux.NewRouter().StrictSlash(true)
	//sample of handling definiton without instrumentation
	//a.router.HandleFunc("/", index).Methods(http.MethodGet)
	//definition of handling with instrumentation
	a.instrumentHandle("/", "index", index).Methods(http.MethodGet)
	a.instrumentHandle("/health", "health", health).Methods(http.MethodGet)
	a.instrumentHandle("/liveness", "liveness", liveness).Methods(http.MethodGet)
	//adding the prometheus metrics
	a.router.Path("/metrics").Handler(promhttp.HandlerFor(a.reg, promhttp.HandlerOpts{}))
	//use handlers for routing
	a.router.Use(mux.CORSMethodMiddleware(a.router))
	log.Debug("Exiting App.init")
}

// index handles the processing of an URL
// w the HTTP writer used to send the response
// r the HTTP request
func index(w http.ResponseWriter, r *http.Request) {
	log.Debug("Entering index(home /) function")
	appName := os.Getenv("APP_NAME")
	version := os.Getenv("APP_VERSION")
	log.Info("processing request on index(/)")
	msg := map[string]string{"message": "Hello from go api server", "application": appName, "version": version}
	utils.RespondJSON(w, r, 200, msg)
	log.Debug("Exiting index(home /) function")
}

// index handles the processing of an URL
// w the HTTP writer used to send the response
// r the HTTP request
func liveness(w http.ResponseWriter, r *http.Request) {
	log.Debug("Entering liveness function")
	log.Info("processing request on index(/liveness)")
	msg := map[string]string{"status": "UP"}
	utils.RespondJSON(w, r, 200, msg)
	log.Debug("Exiting liveness function")

}

// index handles the processing of an URL
// w the HTTP writer used to send the response
// r the HTTP request
func health(w http.ResponseWriter, r *http.Request) {
	log.Debug("Entering health function")
	log.Info("processing request on index(/health)")
	msg := map[string]string{"health": "OK"}
	utils.RespondJSON(w, r, 200, msg)
	log.Debug("Exiting health function")
}

func initLog() {
	myLogLevel := log.InfoLevel
	requestedLevel := os.Getenv("LOG")
	if requestedLevel == "" {
		//do nothing
	} else if requestedLevel == "TRACE" {
		myLogLevel = log.TraceLevel
	} else if requestedLevel == "DEBUG" {
		myLogLevel = log.DebugLevel
	} else if requestedLevel == "INFO" {
		myLogLevel = log.InfoLevel
	} else if requestedLevel == "WARN" {
		myLogLevel = log.WarnLevel
	} else if requestedLevel == "ERROR" {
		myLogLevel = log.ErrorLevel
	} else if requestedLevel == "FATAL" {
		myLogLevel = log.FatalLevel
	} else if requestedLevel == "PANIC" {
		myLogLevel = log.PanicLevel
	}
	log.SetLevel(myLogLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logFile := os.Getenv("LOGFILE")
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		//defer f.Close()
		log.SetOutput(f)
	}
}
