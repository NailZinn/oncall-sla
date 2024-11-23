package main

import (
	"net/http"
	"oncall-sla/metrics"
	"oncall-sla/prober"
	"oncall-sla/sla"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	sla.ConnectToDb()
	sla.InitSliTable()

	go prober.Run()
	go sla.Run()

	reg := prometheus.NewRegistry()

	reg.MustRegister(
		metrics.CreateEventRequestsTotal,
		metrics.CreateEventRequestSuccessTotal,
		metrics.CreateEventRequstFailTotal,
		metrics.CreateEventRequestDurationInMs,
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.ListenAndServe(":2112", nil)
}
