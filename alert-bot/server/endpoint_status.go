package server

import (
	"cosmos-validator-mission-control/alert-bot/config"
	"fmt"
	"log"
	"net/http"
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

	ops = HTTPOptions{
		Endpoint: cfg.LCDEndpoint + "/node_info",
		Method:   http.MethodGet,
	}

	_, err = HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting lcd endpoint status: %v", err)
		msg = msg + fmt.Sprintf("⛔ Unreachable to %s LCD :: and the ERROR is: %v", cfg.LCDEndpoint, err)
	}

	if msg != "" {
		teleErr := SendTelegramAlert(msg, cfg)
		if teleErr != nil {
			log.Printf("Error while sending endpoint status alert to telegram :%v", err)
		}
	}

	return nil

}
