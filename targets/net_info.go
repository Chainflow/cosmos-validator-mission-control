package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"strconv"
)

func GetNetInfo(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}
	var pts []*client.Point

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error getting node_info: %v", err)
		return
	}
	var ni NetInfo
	err = json.Unmarshal(resp.Body, &ni)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	numPeers, err := strconv.Atoi(ni.Result.NumPeers)
	if err != nil {
		log.Printf("Error converting num_peers to int: %v", err)
		numPeers = 0
	} else if int64(numPeers) < cfg.NumPeersThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Number of peers has fallen below %d", cfg.NumPeersThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Number of peers has fallen below %d", cfg.NumPeersThreshold), cfg)
	}
	p1, err := createDataPoint("vcf_num_peers", map[string]string{}, map[string]interface{}{"count": numPeers})
	if err == nil {
		pts = append(pts, p1)
	}

	peerAddrs := make([]string, len(ni.Result.Peers))
	for i, peer := range ni.Result.Peers {
		peerAddrs[i] = peer.RemoteIP
	}
	p2, err := createDataPoint("vcf_peer_addresses", map[string]string{}, map[string]interface{}{"addresses": peerAddrs})
	if err == nil {
		pts = append(pts, p2)
	}

	bp.AddPoints(pts)
	_ = writeBatchPoints(c, bp)
	log.Printf("No. of peers: %d \n Peer Addresses: %v", numPeers, peerAddrs)
}
