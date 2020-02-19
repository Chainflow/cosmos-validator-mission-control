package targets

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GaiadRunningGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "vcf_gaiad_running",
		Help: "Check if gaiad running",
	})
)
