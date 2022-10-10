package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"math/rand"
	"time"
)

var (
	counter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_counts_total",
		Help: "The total number of somethins",
	}, []string{"code", "method"})
	gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "myapp_gauge_total",
		Help:        "Current level of some metric",
		ConstLabels: prometheus.Labels{"session": "design_system"},
	})
	histogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "myapp_random_numbers_hist",
		Help:    "A histogram of normally distributed random numbers.",
		Buckets: prometheus.LinearBuckets(-3, .1, 61), // from -3 with step 0.1
	})
	summary = promauto.NewSummary(prometheus.SummaryOpts{
		Name:       "myapp_random_numbers_summary",
		Help:       "A summary of normally distributed random numbers.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, //quantiles with absolute errors
	})
)

func init() {
	go func() {
		for range time.Tick(100 * time.Millisecond) {
			counter.WithLabelValues("404", "POST").Add(float64(rand.Intn(2)))
			counter.WithLabelValues("404", "GET").Add(float64(rand.Intn(1)))
			counter.WithLabelValues("200", "POST").Add(float64(rand.Intn(10)))
			counter.WithLabelValues("200", "GET").Add(float64(rand.Intn(30)))

			gauge.Set(float64(rand.Intn(10) + 100))
			histogram.Observe(rand.NormFloat64())
			summary.Observe(rand.NormFloat64())
		}
	}()
}
