package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
)

func GetMissedBlocks(cfg *config.Config, c client.Client, cbh int) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(HTTPOptions{
		Endpoint:    cfg.NodeURL + "block",
		QueryParams: QueryParams{"height": strconv.Itoa(cbh)},
		Method:      "GET",
	})
	if err != nil {
		log.Printf("Error getting details of current block: %v", err)
		return
	}

	var b CurrentBlockWithHeight
	err = json.Unmarshal(resp.Body, &b)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	addrExists := false
	for _, c := range b.Result.Block.LastCommit.Precommits {
		if c.ValidatorAddress == cfg.ValidatorAddress {
			addrExists = true
		}
	}

	if !addrExists {
		_ = SendTelegramAlert(fmt.Sprintf("Validator missed a block at block height %d", cbh), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Validator missed a block at block height %d", cbh), cfg)
		_ = writeToInfluxDb(c, bp, "vcf_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh})
	}
}

func GetGaiaCliStatus(_ HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}
	var pts []*client.Point

	cmd := exec.Command("gaiacli", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing cmd gaiacli status")
		return
	}
	var status GaiaCliStatus
	err = json.Unmarshal(out, &status)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var bh int
	currentBlockHeight := status.SyncInfo.LatestBlockHeight
	if currentBlockHeight != "" {
		bh, err = strconv.Atoi(currentBlockHeight)
		p2, err := createDataPoint("vcf_current_block_height", map[string]string{}, map[string]interface{}{"height": bh})
		if err == nil {
			pts = append(pts, p2)
		}
		go GetMissedBlocks(cfg, c, bh)
	}

	var synced int
	caughtUp := !status.SyncInfo.CatchingUp
	if !caughtUp {
		_ = SendTelegramAlert("Your node is not synced!", cfg)
		_ = SendEmailAlert("Your node is not synced!", cfg)
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
