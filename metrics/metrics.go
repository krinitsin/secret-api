package metrics

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var namespace string

var (
	RequestCounter        *prometheus.CounterVec
	RequestDuration       *prometheus.HistogramVec
	SecretApiCall         *prometheus.CounterVec
	SecretApiCallDuration *prometheus.HistogramVec
)

// Init инициализирует метрики прометея
func Init(appName string) {
	namespace = strings.Replace(appName, "-", "_", -1)

	RequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "incoming_request_counter",
		Help:      "Счетчик входящих запросов",
	},
		[]string{"type", "api", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "incoming_request_duration",
		Help:      "Длительность обработки входящего запросов в миллисекундах",
		Buckets:   []float64{0.01, 0.1, 0.5, 1.0, 10},
	},
		[]string{"type", "api", "status"},
	)

	SecretApiCall = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "secret_api_request_counter",
		Help:      "Счетчик запросов к secret api",
	},
		[]string{"type", "status", "resp_code"},
	)

	SecretApiCallDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "secret_api_request_duration",
		Help:      "Длительность запросов к secret api в миллисекундах",
		Buckets:   []float64{0.01, 0.1, 0.5, 1.0, 10},
	},
		[]string{"type", "status", "resp_code"},
	)

	prometheus.MustRegister(
		RequestCounter,
		RequestDuration,
		SecretApiCall,
		SecretApiCallDuration,
	)
}

// AddRequest добавляет входящий запрос
func AddRequest(method string, api string, status string, duration time.Duration) {
	RequestCounter.WithLabelValues(method, api, status).Inc()
	RequestDuration.WithLabelValues(method, api, status).Observe(duration.Seconds() * 1000)
}

// AddSecretApiCall добавляет завтрос к внешнему API secret
func AddSecretApiCall(method, status string, code int, duration time.Duration) {
	SecretApiCall.WithLabelValues(method, status,string(code)).Inc()
	SecretApiCallDuration.WithLabelValues(method, status, string(code)).Observe(duration.Seconds() * 1000)
}
