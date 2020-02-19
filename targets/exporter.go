package targets

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GaiadRunningGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "vcf_gaiad_running",
		Help: "Check if gaiad running",
	})
	NumPeersGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "vcf_num_peers",
		Help: "Check number of peers",
	})
	PeerAddressesSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "vcf_peer_addresses",
		Help: "Check peer addresses of node",
	}, []string{"ips"})
)
