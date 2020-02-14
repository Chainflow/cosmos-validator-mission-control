package monitor

import (
	"encoding/json"
	"log"
)

func Address(m HTTPOptions) {
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

	log.Printf("Address: %s", validatorResp.Validator.Uptime.Address)
}
