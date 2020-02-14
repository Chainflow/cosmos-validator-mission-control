package monitor

import (
	"encoding/json"
	"log"
	"strconv"
)

func MaxChangeRate(m HTTPOptions) {
	resp, err := RunMonitor(m)
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

	maxChangeRate, err := strconv.Atoi(validatorResp.Validator.Details.Commission.MaxChangeRate)
	if err != nil {
		log.Printf("Error converting max change rate %s to percentage", validatorResp.Validator.Details.Commission.MaxChangeRate)
		return
	}

	log.Printf("Max Change Rate: %s%", maxChangeRate*100)
}
