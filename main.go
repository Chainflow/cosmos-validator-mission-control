package main

import (
	"chainflow-vitwit/config"
	"chainflow-vitwit/targets"
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"sync"
	"time"
)

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

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http:// localhost:%s", cfg.InfluxDB.Port),
		Username: cfg.InfluxDB.Username,
		Password: cfg.InfluxDB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	var wg sync.WaitGroup
	for _, tg := range m.List {
		wg.Add(1)
		go func(target targets.Target) {
			for {
				runner.Run(target.Func, target.HTTPOptions, cfg, c)
				time.Sleep(scrapeRate)
			}
		}(tg)
	}
	wg.Wait()
}
