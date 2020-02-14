package monitor

import (
	"encoding/json"
	"log"
)

func PeerAddresses(m HTTPOptions) {
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

	pAddrs := make([]string, len(ni.Result.Peers))
	for i, peer := range ni.Result.Peers {
		pAddrs[i] = peer.RemoteIP
	}

	log.Printf("Peer Addresses: %v", pAddrs)
}
