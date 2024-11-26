package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequestCountWithPath = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_with_path",
			Help: "Number of HTTP requests by path.",
		},
		[]string{"path"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Average response time of HTTP requests.",
		},
		[]string{"path"},
	)

	KafkaMessagesReceivedCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_received_total",
			Help: "Total number of messages received from Kafka.",
		},
		[]string{"topic"},
	)

	KafkaMessageProcessingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "kafka_message_processing_duration_seconds",
			Help: "Duration of processing Kafka messages.",
		},
		[]string{"topic"},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestCountWithPath)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(KafkaMessagesReceivedCount)
	prometheus.MustRegister(KafkaMessageProcessingDuration)
}
