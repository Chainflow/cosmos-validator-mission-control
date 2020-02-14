package monitor

import (
	"encoding/json"
	"log"
)

func AddressBalance(m HTTPOptions) {
	resp, err := RunMonitor(m)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var accResp AccountResp
	err = json.Unmarshal(resp.Body, &accResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	log.Printf("Address Balance: %s", accResp.Account.Balance[0].Amount+accResp.Account.Balance[0].Denom)
}
