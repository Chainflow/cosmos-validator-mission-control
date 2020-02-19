package main

import (
	"chainflow-vitwit/config"
	"chainflow-vitwit/targets"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func init() {
	prometheus.MustRegister(targets.GaiadRunningGauge)
	prometheus.MustRegister(targets.NumPeersGauge)
}

func main() {
	cfg, err := config.ReadFromFile()
	if err != nil {
		log.Fatal(err)
	}

	m := targets.InitTargets(cfg)
	runner := targets.NewRunner()
	scrapeRate, err := time.ParseDuration(cfg.Scraper.Rate)
	if err != nil {
		log.Fatal(err)
	}

	for _, tg := range m.List {
		go func(target targets.Target) {
			for {
				runner.Run(target.Func, target.HTTPOptions, cfg)
				time.Sleep(scrapeRate)
			}
		}(tg)
	}

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(cfg.Scraper.Port, nil))
}
