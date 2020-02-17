package targets

import (
	"encoding/json"
	"log"
)

func GetNetInfo(ops HTTPOptions) {
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	var ni NetInfo
	err = json.Unmarshal(resp.Body, &ni)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	numPeers := ni.Result.NumPeers

	peerAddrs := make([]string, len(ni.Result.Peers))
	for i, peer := range ni.Result.Peers {
		peerAddrs[i] = peer.RemoteIP
	}

	log.Printf("No. of peers: %s \n Peer Addresses: %v", numPeers, peerAddrs)
}
