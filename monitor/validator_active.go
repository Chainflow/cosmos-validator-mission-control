package monitor

import (
	"encoding/json"
	"log"
	"os/exec"
)

func IsValidatorActive(_ HTTPOptions) {
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

	if status.NodeInfo.ValidatorInfo.VotingPower == "0" {
		log.Println("Validator is INACTIVE and Jailed!")
		return
	}

	log.Println("Validator is active")
}
