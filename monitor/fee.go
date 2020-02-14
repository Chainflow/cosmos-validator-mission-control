package monitor

import (
	"encoding/json"
	"log"
)

func Fee(m HTTPOptions) {
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

	log.Printf("Fee: %s", validatorResp.Validator.Details.Commission.Rate)
}
