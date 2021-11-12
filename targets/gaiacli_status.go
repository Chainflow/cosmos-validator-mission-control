package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetGaiaCliStatus to reponse of validator status like
//current block height and node status
func GetGaiaCliStatus(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}
	var pts []*client.Point

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var status ValidatorRpcStatus
	err = json.Unmarshal(resp.Body, &status)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var bh int
	currentBlockHeight := status.Result.SyncInfo.LatestBlockHeight
	if currentBlockHeight != "" {
		bh, _ = strconv.Atoi(currentBlockHeight)
		p2, err := createDataPoint("vcf_current_block_height", map[string]string{}, map[string]interface{}{"height": bh})
		if err == nil {
			pts = append(pts, p2)
		}
	}

	var synced int
	caughtUp := !status.Result.SyncInfo.CatchingUp
	if !caughtUp {
		_ = SendTelegramAlert(fmt.Sprintf("%s Your validator node is not synced!", cfg.ValidatorName), cfg)
		_ = SendEmailAlert(fmt.Sprintf("%s Your validator node is not synced!", cfg.ValidatorName), cfg)
		synced = 0
	} else {
		synced = 1
	}
	p3, err := createDataPoint("vcf_node_synced", map[string]string{}, map[string]interface{}{"status": synced})
	if err == nil {
		pts = append(pts, p3)
	}

	bp.AddPoints(pts)
	_ = writeBatchPoints(c, bp)
	log.Printf("\nCurrent Block Height: %s \nCaught Up? %t \n",
		currentBlockHeight, caughtUp)
}
