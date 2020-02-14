package main

import (
	"chainflow-vitwit/config"
	"chainflow-vitwit/monitor"
	"log"
	"sync"
	"time"
)

func main() {
	cfg, err := config.ReadFromTomlFile()
	if err != nil {
		log.Fatal(err)
	}
	m := monitor.InitNodeMonitors(cfg)
	var wg sync.WaitGroup

	for _, monit := range m.List {
		wg.Add(1)
		t := time.Tick(5 * time.Second)
		go func(t <-chan time.Time, m monitor.NodeMonitor) {
			for _ = range t {
				m.Func(m.HTTPOptions)
			}
		}(t, monit)
	}
	wg.Wait()
}
