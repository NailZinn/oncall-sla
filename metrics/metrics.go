package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OncallCreateEventRequestsTotal        = "oncall_create_event_requests_total"
	OncallCreateEventRequestsSuccessTotal = "oncall_create_event_requests_success_total"
	OncallCreateEventRequestsFailTotal    = "oncall_create_event_requests_fail_total"
	OncallCreateEventRequestsDurationMs   = "oncall_create_event_request_duration_ms"

	CreateEventRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: OncallCreateEventRequestsTotal,
		Help: "The total number of event creation requests",
	})
	CreateEventRequestSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: OncallCreateEventRequestsSuccessTotal,
		Help: "The total number of successful event creation requests",
	})
	CreateEventRequstFailTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: OncallCreateEventRequestsFailTotal,
		Help: "The total number of failed event creation requests",
	})
	CreateEventRequestDurationInMs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: OncallCreateEventRequestsDurationMs,
		Help: "Duration in ms of event creation request",
	})
)
