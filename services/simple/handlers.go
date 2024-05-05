package simple

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/massimo-gollo/benchy/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"strconv"
	"time"
)

var (
	computationTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "benchy_computation_time",
			Help: "Time taken to execute the Fibonacci function.",
		},
		[]string{"function", "endpoint", "podId"},
	)

	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "benchy_request",
			Help: "Total number of requests received",
		},
		[]string{"method", "endpoint", "podId"},
	)

	queryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "benchy_query_duration",
			Help: "Benchy query durations in nanoseconds",
			Buckets: []float64{
				0.00001,
				0.00002,
				0.00003,
				0.00004,
				0.00005,
				0.00006,
				0.00007,
				0.00008,
				0.00009,
				0.0001,
				0.0002,
				0.0003,
				0.0004,
				0.0005,
				0.0006,
				0.0007,
				0.0008,
				0.0009,
				0.001,
				0.002,
				0.003,
				0.004,
				0.005,
				0.006,
				0.007,
				0.008,
				0.009,
				0.01,
				0.02,
				0.03,
				0.04,
				0.05,
				0.06,
				0.07,
				0.08,
				0.09,
				0.1,
				0.2,
				0.3,
				0.4,
				0.5,
				0.6,
				0.7,
				0.8,
				0.9,
				1.0,
			},
		},
		[]string{"endpoint", "podId"},
	)
)

func init() {
	prometheus.MustRegister(computationTime)
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(queryDuration)
}

func HealthHandler(c *gin.Context) {
	defer func() {
		envBenchyPodId := os.Getenv("BenchyPodID")
		requestCounter.WithLabelValues("GET", "/healthz", envBenchyPodId).Inc()
	}()
	c.JSON(200, gin.H{
		"health": "ok",
	})
}

func PrometheusHandler(c *gin.Context) {

	registry := prometheus.NewRegistry()
	registry.MustRegister(computationTime)
	registry.MustRegister(requestCounter)
	registry.MustRegister(queryDuration)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(c.Writer, c.Request)

}
func CpuTaskHandler(c *gin.Context) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		envBenchyPodId := os.Getenv("BenchyPodID")
		computationTime.WithLabelValues("FibonacciTime", "/cpu", envBenchyPodId).Set(elapsed.Seconds())
		requestCounter.WithLabelValues("GET", "/cpu", envBenchyPodId).Inc()
		queryDuration.WithLabelValues("/cpu", envBenchyPodId).Observe(float64(elapsed.Seconds()))
	}()
	n, _ := strconv.Atoi(c.Query("n"))
	result := util.FibonacciOptimized(uint64(n))
	elapsed := time.Since(start)
	c.JSON(200, gin.H{
		"result": fmt.Sprintf("Fibonacci sequence of %d: %d\n. Elapsed: %s", n, result, elapsed),
	})
}

func CpuIntensiveTaskHandler(c *gin.Context) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		envBenchyPodId := os.Getenv("BenchyPodID")
		computationTime.WithLabelValues("FibonacciTime_intensive", "/cpuintensive", envBenchyPodId).Set(elapsed.Seconds())
		requestCounter.WithLabelValues("GET", "/cpuintensive", envBenchyPodId).Inc()
		queryDuration.WithLabelValues("/cpuintensive", envBenchyPodId).Observe(float64(elapsed.Seconds()))
	}()
	n, _ := strconv.Atoi(c.Query("n"))
	result := util.Fibonacci(n)
	elapsed := time.Since(start)
	c.JSON(200, gin.H{
		"result": fmt.Sprintf("Fibonacci sequence of %d: %d\n. Elapsed: %s", n, result, elapsed),
	})
}

func MemTaskHandler(c *gin.Context) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		envBenchyPodId := os.Getenv("BenchyPodID")
		computationTime.WithLabelValues("MemTime", "/mem", envBenchyPodId).Set(elapsed.Seconds())
		requestCounter.WithLabelValues("GET", "/mem", envBenchyPodId).Inc()
		queryDuration.WithLabelValues("/mem", envBenchyPodId).Observe(float64(elapsed.Seconds()))
	}()
	size, _ := strconv.Atoi(c.Query("size"))
	data := util.GenerateData(size)
	c.JSON(200, gin.H{
		"result": fmt.Sprintf("Generated %d random numbers.\n", len(data)),
	})
}
