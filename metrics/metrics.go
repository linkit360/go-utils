package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// name of the service / app to record them in prometheus
var InstancePrefix string

func Init(instancePrefix string) {
	if instancePrefix == "" {
		panic("instance prefix is empty")
	}
	InstancePrefix = instancePrefix
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
func NewGauge(namespace, subsystem, name, help string) Gauge {
	g := Gauge{}
	g.gauge = PrometheusGauge(namespace, subsystem, name, help)
	return g
}
func NewGaugeLabel(namespace, subsystem, name, help string, labels map[string]string) Gauge {
	g := Gauge{}
	g.gauge = PrometheusGaugeLabel(namespace, subsystem, name, help, labels)
	return g
}
func PrometheusGauge(namespace, subsystem, name, help string) prometheus.Gauge {
	if InstancePrefix == "" {
		panic("instance prefix is empty")
	}
	if namespace == "" {
		namespace = InstancePrefix
	} else {
		namespace = InstancePrefix + "_" + namespace
	}
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
	if InstancePrefix == "" {
		panic("instance prefix is empty")
	}
	if namespace == "" {
		namespace = InstancePrefix
	} else {
		namespace = InstancePrefix + "_" + namespace
	}
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
	if InstancePrefix == "" {
		panic("instance prefix is empty")
	}
	if name == "" {
		name = InstancePrefix
	} else {
		name = InstancePrefix + "_" + name
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
