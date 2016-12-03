package metrics

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// name of the service / app to record them in prometheus
var appName string

func Init(name string) {
	appName = name
}

func AddHandler(r *gin.Engine) {
	rg := r.Group("/metrics")
	rg.GET("", gin.WrapH(prometheus.Handler()))
}

type Gauge struct {
	gauge   prometheus.Gauge
	counter int64
}

func (g *Gauge) Inc() {
	g.counter++
}
func (g *Gauge) Update() {
	g.gauge.Set(float64(g.counter))
	g.counter = 0
}

// for any gauges
func NewGaugeAlert(namespace, subsystem, name, help string) Gauge {
	g := Gauge{}
	g.gauge = PrometheusGauge(namespace, subsystem, name, help)
	go func() {
		for range time.Tick(time.Minute) {
			g.Update()
		}
	}()
	return g
}
func NewGauge(namespace, subsystem, name, help string) Gauge {
	return NewGaugeAlert(namespace, subsystem, name, help)
}

func PrometheusGauge(namespace, subsystem, name, help string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})
	prometheus.MustRegister(gauge)
	return gauge
}
func PrometheusGaugeLabel(namespace, subsystem, name, help string, labels map[string]string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: labels,
	})
	prometheus.MustRegister(gauge)
	return gauge
}

// for duration
func NewSummary(name, help string) prometheus.Summary {
	if appName == "" {
		log.Fatal("app name is empty")
	}
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: name,
			Help: help,
		},
	)
	prometheus.MustRegister(summary)
	return summary
}
