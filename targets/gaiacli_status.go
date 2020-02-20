package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"os/exec"
	"strconv"
)

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

	var p1 *client.Point
	validatorActive := status.NodeInfo.ValidatorInfo.VotingPower != "0"
	if !validatorActive {
		_ = SendTelegramAlert("Validator has been jailed!", cfg)
		_ = SendEmailAlert("Validator has been jailed!", cfg)
		p1, err = createDataPoint("vcf_validator_status", map[string]string{}, map[string]interface{}{"status": 0})
	} else {
		p1, err = createDataPoint("vcf_validator_status", map[string]string{}, map[string]interface{}{"status": 1})
	}
	if err == nil {
		pts = append(pts, p1)
	}

	var bh int
	currentBlockHeight := status.NodeInfo.SyncInfo.LatestBlockHeight
	if currentBlockHeight == "" {
		bh = 0
	} else {
		bh, err = strconv.Atoi(currentBlockHeight)
		if err != nil {
			bh = 0
		}
	}
	p2, err := createDataPoint("vcf_current_block_height", map[string]string{}, map[string]interface{}{"height": bh})
	if err == nil {
		pts = append(pts, p2)
	}

	var synced int
	caughtUp := !status.NodeInfo.SyncInfo.CatchingUp
	if caughtUp {
		_ = SendTelegramAlert("Your node has been synced!", cfg)
		_ = SendEmailAlert("Your node has been synced!", cfg)
		synced = 1
	} else {
		synced = 0
	}
	p3, err := createDataPoint("vcf_node_synced", map[string]string{}, map[string]interface{}{"status": synced})
	if err == nil {
		pts = append(pts, p3)
	}

	var vp int
	if status.NodeInfo.ValidatorInfo.VotingPower != "" {
		vp, err = strconv.Atoi(status.NodeInfo.ValidatorInfo.VotingPower)
		if err != nil {
			log.Printf("Error while converting votingPower to int: %v", err)
			vp = 0
		}
	} else {
		vp = 0
	}
	p4, err := createDataPoint("vcf_voting_power", map[string]string{}, map[string]interface{}{"power": vp})
	if err == nil {
		pts = append(pts, p4)
	}
	if int64(vp) <= cfg.VotingPowerThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Your voting power has dropped below %d", cfg.VotingPowerThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Your voting power has dropped below %d", cfg.VotingPowerThreshold), cfg)
	}

	bp.AddPoints(pts)
	_ = writeBatchPoints(c, bp)
	log.Printf("Validator Active: %t \nCurrent Block Height: %s \nCaught Up? %t \nVoting Power: %d \n",
		validatorActive, currentBlockHeight, caughtUp, vp)
}
