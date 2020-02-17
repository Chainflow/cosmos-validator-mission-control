package targets

import (
	"encoding/json"
	"log"
	"os/exec"
)

func GetGaiaCliStatus(_ HTTPOptions) {
	cmd := exec.Command("gaiacli", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		return
	}

	var status GaiaCliStatus
	err = json.Unmarshal(out, &status)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	validatorActive := status.NodeInfo.ValidatorInfo.VotingPower != "0"

	currentBlockHeight := status.NodeInfo.SyncInfo.LatestBlockHeight

	caughtUp := !status.NodeInfo.SyncInfo.CatchingUp

	votingPower := status.NodeInfo.ValidatorInfo.VotingPower

	log.Printf("Validator Active: %t \nCurrent Block Height: %s \nCaught Up? %t \nVoting Power: %s \n",
		validatorActive, currentBlockHeight, caughtUp, votingPower)
}
