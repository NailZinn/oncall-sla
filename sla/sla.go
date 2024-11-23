package sla

import (
	"fmt"
	"oncall-sla/metrics"
	"time"
)

func Run() {
	for {
		time.Sleep(1 * time.Minute)
		now := int32(time.Now().UTC().UnixMilli() / 1000)

		query := fmt.Sprintf("changes(%s[1m])", metrics.OncallCreateEventRequestsTotal)
		metric := int32(fetchMetrics(query, now, 0))
		saveSli(metrics.OncallCreateEventRequestsTotal, float32(metric), 1, now, metric < 1)

		query = fmt.Sprintf("changes(%s[1m])", metrics.OncallCreateEventRequestsSuccessTotal)
		metric = int32(fetchMetrics(query, now, 0))
		saveSli(metrics.OncallCreateEventRequestsSuccessTotal, float32(metric), 1, now, metric < 1)

		query = fmt.Sprintf("changes(%s[1m])", metrics.OncallCreateEventRequestsFailTotal)
		metric = int32(fetchMetrics(query, now, 1))
		saveSli(metrics.OncallCreateEventRequestsFailTotal, float32(metric), 0, now, metric > 0)

		metric = int32(fetchMetrics(metrics.OncallCreateEventRequestsDurationMs, now, 200))
		saveSli(metrics.OncallCreateEventRequestsDurationMs, float32(metric), 50, now, metric > 50)
	}
}
