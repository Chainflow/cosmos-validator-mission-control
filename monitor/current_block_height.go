package monitor

import (
	"encoding/json"
	"log"
	"os/exec"
)

func CurrentBlockHeight(_ HTTPOptions) {
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
	log.Printf("Current Block Height: %s", status.NodeInfo.SyncInfo.LatestBlockHeight)
}
