package server

import (
	"cosmos-validator-mission-control/alert-bot/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// JailedAlerting to send transaction alert to telegram and mail
func ValidatorStatusAlert(ops HTTPOptions, cfg *config.Config, c client.Client) {
	log.Println("Coming inside validator status alerting")
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var validatorResp ValidatorResp
	err = json.Unmarshal(resp.Body, &validatorResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
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

	validatorStatus := validatorResp.Result.Jailed

	if !validatorStatus {
		_ = writeToInfluxDb(c, bp, "vcf_val_status", map[string]string{}, map[string]interface{}{"status": "voting"})
		if t == a1 || t == a2 {
			_ = SendTelegramAlert(fmt.Sprintf("Your validator %s is currently voting", cfg.ValidatorName), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Your validator is currently voting"), cfg)
			log.Println("Sent validator status alert")
		}
	} else {
		_ = writeToInfluxDb(c, bp, "vcf_val_status", map[string]string{}, map[string]interface{}{"status": "jailed"})
		if t == a1 || t == a2 {
			_ = SendTelegramAlert(fmt.Sprintf("Your validator %s is in jailed status", cfg.ValidatorName), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Your validator is in jailed status"), cfg)
			log.Println("Sent validator status alert")
		}
	}
	return
}

// CheckValidatorJailed to send transaction alert to telegram and mail
// when the validator will be jailed
func CheckValidatorJailed(cfg *config.Config) error {
	log.Println("Coming inside jailed alerting")
	ops := HTTPOptions{
		Endpoint: cfg.LCDEndpoint + "/staking/validators/" + cfg.ValOperatorAddress,
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
		_ = SendTelegramAlert(fmt.Sprintf("%s validator is in jailed status", cfg.ValidatorName), cfg)
		_ = SendEmailAlert(fmt.Sprintf("%s validator is in jailed status", cfg.ValidatorName), cfg)
		log.Println("Sent validator jailed status alert")
	}
	return nil
}

// GetValStatusFromDB returns latest current height from db
func GetValStatusFromDB(cfg *config.Config, c client.Client) string {
	var valStatus string
	q := client.NewQuery("SELECT last(status) FROM vcf_val_status", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						status := r.Series[0].Values[0][idx]
						valStatus = fmt.Sprintf("%v", status)
						break
					}
				}
			}
		}
	}
	return valStatus
}
