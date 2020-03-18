package server

import (
	"chainflow-vitwit/alert-bot/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Send transaction alert to telegram and mail
func JailedAlerting(cfg *config.Config) error {
	log.Println("Coming inside jailed alerting")
	ops := HTTPOptions{
		Endpoint: cfg.LCDEndpoint + "staking/validators/" + cfg.OperatorAddress,
		Method:   http.MethodGet,
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	var validatorResp ValidatorResp
	err = json.Unmarshal(resp.Body, &validatorResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	validatorStatus := validatorResp.Result.Jailed
	if validatorStatus == false {
		_ = SendTelegramAlert(fmt.Sprintf("Your validator is in active status"), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Your validator is in active status"), cfg)
		log.Println("Sent validator status alert")
	} else {
		_ = SendTelegramAlert(fmt.Sprintf("Your validator is in jailed status"), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Your validator is in jailed status"), cfg)
		log.Println("Sent validator status alert")
	}
	return nil
}
