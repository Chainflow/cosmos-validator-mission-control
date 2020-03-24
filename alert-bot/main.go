package main

import (
	"chainflow-vitwit/alert-bot/config"
	"chainflow-vitwit/alert-bot/server"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	cfg, err := config.ReadFromFile()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			if err := server.SendSingleMissedBlockAlert(cfg); err != nil {
				fmt.Println("Error while sending missed block alerts", err)
			}
			time.Sleep(4 * time.Second)
		}
	}()

	go func() {
		for {
			if err := server.ValidatorStatusAlert(cfg); err != nil {
				fmt.Println("Error while sending jailed alerts", err)
			}
			time.Sleep(60 * time.Second)
		}
	}()

	wg.Wait()
}
