package monitor

import (
	"encoding/json"
	"log"
)

func OperatorAddress(m HTTPOptions) {
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

	log.Printf("Operator Address: %s", validatorResp.Validator.Details.OperatorAddress)
}
