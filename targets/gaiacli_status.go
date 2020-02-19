package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func GetGaiaCliStatus(_ HTTPOptions, cfg *config.Config) {
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

	currentBlockHeight := status.NodeInfo.SyncInfo.LatestBlockHeight

	caughtUp := !status.NodeInfo.SyncInfo.CatchingUp
	if caughtUp {
		_ = SendTelegramAlert("Your node has been synced!", cfg)
		_ = SendEmailAlert("Your node has been synced!", cfg)
	}

	vp, err := strconv.Atoi(status.NodeInfo.ValidatorInfo.VotingPower)
	if err != nil {
		log.Printf("Error while converting votingPower to int: %v", err)
	}
	if int64(vp) <= cfg.VotingPowerThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Your voting power has dropped below %d", cfg.VotingPowerThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Your voting power has dropped below %d", cfg.VotingPowerThreshold), cfg)
	}

	log.Printf("Validator Active: %t \nCurrent Block Height: %s \nCaught Up? %t \nVoting Power: %d \n",
		validatorActive, currentBlockHeight, caughtUp, vp)
}
