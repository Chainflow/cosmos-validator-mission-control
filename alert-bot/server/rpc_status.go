package server

import (
	"cosmos-validator-mission-control/alert-bot/config"
	"fmt"
	"log"
)

func GetEndpointStatus(cfg *config.Config) error {
	var msg string

	ops := HTTPOptions{
		Endpoint: cfg.ExternalRPC + "/status",
		Method:   "GET",
	}

	_, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		msg = msg + fmt.Sprintf("⛔ Unreachable to %s RPC :: and the ERROR is: %v", cfg.ExternalRPC, err)
	}

	if msg != "" {
		teleErr := SendTelegramAlert(msg, cfg)
		if teleErr != nil {
			log.Printf("Error while sending endpoint status alert to telegram :%v", err)
		}
	}

	return nil

}
