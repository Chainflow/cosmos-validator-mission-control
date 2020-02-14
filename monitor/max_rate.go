package monitor

import (
	"encoding/json"
	"log"
	"strconv"
)

func MaxRate(m HTTPOptions) {
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

	maxRate, err := strconv.Atoi(validatorResp.Validator.Details.Commission.MaxRate)
	if err != nil {
		log.Printf("Error converting max rate %s to percentage", validatorResp.Validator.Details.Commission.MaxRate)
		return
	}

	log.Printf("Max Rate: %s%", maxRate*100)
}
