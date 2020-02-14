package monitor

import (
	"encoding/json"
	"log"
)

func NumPeers(m HTTPOptions) {
	resp, err := RunMonitor(m)
	if err != nil {
		// send alert
		log.Printf("Error: %v", err)
		return
	}
	var ni NetInfo
	err = json.Unmarshal(resp.Body, &ni)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	log.Printf("Num Peers: %s", ni.Result.NumPeers)
}
