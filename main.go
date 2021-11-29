package main

import (
	"flag"
	"fmt"
	"math/rand" // #nosec
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"nginx_exporter/metric"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

var (
	/*	EnableMetrics, _ = strconv.ParseBool(os.Getenv("EnableMetrics"))
		MetricsPerHost, _ = strconv.ParseBool(os.Getenv("MetricsPerHost"))
		ListenPorts, _ = strconv.ParseInt(os.Getenv("ListenPorts"), 10, 64)*/

	EnableMetrics  = true
	MetricsPerHost = true
	ListenPorts    int
	err            error
	// 新增  http_stub_status_module 模块的path和port
	NginxStatusPath string
	NginxStatusPort string
)

func main() {

	klog.InitFlags(nil)
	//flag.Set("logtostderr", "false")     // By default klog logs to stderr, switch that off
	//flag.Set("alsologtostderr", "false") // false is default, but this is informative
	//flag.Set("stderrthreshold", "FATAL") // stderrthreshold defaults to ERROR, we don't want anything in stderr
	//flag.Set("log_file", "myfile.log")   // log to a file

	// http 启动端口，默认9123
	flag.IntVar(&ListenPorts, "port", 9123, "http监听端口,默认9123")

	// http_stub_status_module 模块的path和port
	flag.StringVar(&NginxStatusPath, "statuspath", "/stub_status", "http_stub_status_module 模块的监听路径,默认/stub_status")
	flag.StringVar(&NginxStatusPort, "statusport", "8021", "http_stub_status_module 模块的监听端口,默认8021")

	// parse klog/v2 flags
	flag.Parse()

	// make sure we flush before exiting
	defer klog.Flush()

	rand.Seed(time.Now().UnixNano())

	reg := prometheus.NewRegistry()

	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{
		PidFn:        func() (int, error) { return os.Getpid(), nil },
		ReportErrors: true,
	}))

	mc := metric.NewDummyCollector()

	if EnableMetrics {
		mc, err = metric.NewCollector(NginxStatusPath, NginxStatusPort, MetricsPerHost, reg)
		if err != nil {
			println(time.Now().Format(time.UnixDate), ": ", "Error creating prometheus collector:  %v", err)
		}
	}
	mc.Start()

	mux := http.NewServeMux()
	registerMetrics(reg, mux)
	registerProfiler(mux)

	go startHTTPServer(int(ListenPorts), mux)

	for {
		time.Sleep(time.Second * 60)
		//Println(time.Now().Format(time.UnixDate), "the NumGoroutine done is: ", runtime.NumGoroutine())
	}
}

func registerMetrics(reg *prometheus.Registry, mux *http.ServeMux) {
	mux.Handle(
		"/metrics",
		promhttp.InstrumentMetricHandler(
			reg,
			promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		),
	)
}

func registerProfiler(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/heap", pprof.Index)
	mux.HandleFunc("/debug/pprof/mutex", pprof.Index)
	mux.HandleFunc("/debug/pprof/goroutine", pprof.Index)
	mux.HandleFunc("/debug/pprof/threadcreate", pprof.Index)
	mux.HandleFunc("/debug/pprof/block", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func startHTTPServer(port int, mux *http.ServeMux) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	println(time.Now().Format(time.UnixDate), ": ", server.ListenAndServe())
}
