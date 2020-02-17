package main

import (
	"chainflow-vitwit/config"
	"chainflow-vitwit/targets"
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
	var wg sync.WaitGroup

	for _, tg := range m.List {
		wg.Add(1)
		t := time.Tick(5 * time.Second)
		go func(t <-chan time.Time, target targets.Target) {
			for _ = range t {
				target.Func(target.HTTPOptions, cfg)
			}
		}(t, tg)
	}
	wg.Wait()
}
