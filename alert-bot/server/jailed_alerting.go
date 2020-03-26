package server

import (
	"chainflow-vitwit/alert-bot/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// JailedAlerting to send transaction alert to telegram and mail
func ValidatorStatusAlert(cfg *config.Config) error {
	log.Println("Coming inside validator status alerting")
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

	alertTime1 := cfg.AlertTime1
	alertTime2 := cfg.AlertTime2

	t1, _ := time.Parse(time.Kitchen, alertTime1)
	t2, _ := time.Parse(time.Kitchen, alertTime2)

	now := time.Now().UTC()
	t := now.Format(time.Kitchen)

	a1 := t1.Format(time.Kitchen)
	a2 := t2.Format(time.Kitchen)

	log.Println("a1, a2 and present time : ", a1, a2, t)

	if t == a1 || t == a2 {
		validatorStatus := validatorResp.Result.Jailed
		if !validatorStatus {
			_ = SendTelegramAlert(fmt.Sprintf("Your validator is currently voting"), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Your validator is currently voting"), cfg)
			log.Println("Sent validator status alert")
		} else {
			_ = SendTelegramAlert(fmt.Sprintf("Your validator is in jailed status"), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Your validator is in jailed status"), cfg)
			log.Println("Sent validator status alert")
		}
	}
	return nil
}

// JailedTxAlerting to send transaction alert to telegram and mail
// when the validator will be jailed
func JailedTxAlerting(cfg *config.Config) error {
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
	if validatorStatus {
		_ = SendTelegramAlert(fmt.Sprintf("Your validator is in jailed status"), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Your validator is in jailed status"), cfg)
		log.Println("Sent validator jailed status alert")
	}
	return nil
}
