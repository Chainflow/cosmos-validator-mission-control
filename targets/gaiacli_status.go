package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"log"
	"os/exec"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetGaiaCliStatus to run command gaiacli status and handle the reponse of it
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
	}

	var synced int
	caughtUp := !status.SyncInfo.CatchingUp
	if !caughtUp {
		_ = SendTelegramAlert("Your validator node is not synced!", cfg)
		_ = SendEmailAlert("Your validator node is not synced!", cfg)
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
