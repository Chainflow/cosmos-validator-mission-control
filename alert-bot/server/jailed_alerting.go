package server

import (
	"cosmos-validator-mission-control/alert-bot/config"
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
		Endpoint: cfg.LCDEndpoint + "/cosmos/staking/v1beta1/validators/" + cfg.ValOperatorAddress,
		Method:   http.MethodGet,
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	var validatorResp Validator
	err = json.Unmarshal(resp.Body, &validatorResp)
	if err != nil {
		log.Printf("Error while unmarshelling staking val res : %v", err)
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
		validatorStatus := validatorResp.Validator.Jailed
		if !validatorStatus {
			_ = SendTelegramAlert(fmt.Sprintf("%s validator is currently voting", cfg.ValidatorName), cfg)
			_ = SendEmailAlert(fmt.Sprintf("%s validator is currently voting", cfg.ValidatorName), cfg)
			log.Println("Sent validator status alert")
		} else {
			_ = SendTelegramAlert(fmt.Sprintf("%s validator is in jailed status", cfg.ValidatorName), cfg)
			_ = SendEmailAlert(fmt.Sprintf("%s validator is in jailed status", cfg.ValidatorName), cfg)
			log.Println("Sent validator status alert")
		}
	}
	return nil
}

// CheckValidatorJailed to send transaction alert to telegram and mail
// when the validator will be jailed
func CheckValidatorJailed(cfg *config.Config) error {
	log.Println("Coming inside jailed alerting")
	ops := HTTPOptions{
		Endpoint: cfg.LCDEndpoint + "/cosmos/staking/v1beta1/validators/" + cfg.ValOperatorAddress,
		Method:   http.MethodGet,
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	var validatorResp Validator
	err = json.Unmarshal(resp.Body, &validatorResp)
	if err != nil {
		log.Printf("Error while unmarshelling staking val res : %v", err)
		return err
	}

	validatorStatus := validatorResp.Validator.Jailed
	if validatorStatus {
		_ = SendTelegramAlert(fmt.Sprintf("%s validator is in jailed status", cfg.ValidatorName), cfg)
		_ = SendEmailAlert(fmt.Sprintf("%s validator is in jailed status", cfg.ValidatorName), cfg)
		log.Println("Sent validator jailed status alert")
	}
	return nil
}
