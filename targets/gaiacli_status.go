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

	validatorActive := status.NodeInfo.ValidatorInfo.VotingPower != "0"
	if !validatorActive {
		_ = SendTelegramAlert("Validator has been jailed!", cfg)
		_ = SendEmailAlert("Validator has been jailed!", cfg)
	}
	p1, err := createDataPoint("vcf_validator_status", map[string]string{}, map[string]interface{}{"status": validatorActive})
	if err == nil {
		pts = append(pts, p1)
	}

	currentBlockHeight := status.NodeInfo.SyncInfo.LatestBlockHeight
	if currentBlockHeight != "" {
		p2, err := createDataPoint("vcf_current_block_height", map[string]string{}, map[string]interface{}{"height": currentBlockHeight})
		if err == nil {
			pts = append(pts, p2)
		}
	}

	caughtUp := !status.NodeInfo.SyncInfo.CatchingUp
	if caughtUp {
		_ = SendTelegramAlert("Your node has been synced!", cfg)
		_ = SendEmailAlert("Your node has been synced!", cfg)
	}
	p3, err := createDataPoint("vcf_node_synced", map[string]string{}, map[string]interface{}{"status": caughtUp})
	if err == nil {
		pts = append(pts, p3)
	}

	vp, err := strconv.Atoi(status.NodeInfo.ValidatorInfo.VotingPower)
	if err != nil {
		log.Printf("Error while converting votingPower to int: %v", err)
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
