package monitor

import (
	"log"
	"os/exec"
	"regexp"
)

func GaiadVersion(_ HTTPOptions) {
	cmd := exec.Command("gaiad", "version", "--long")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		return
	}

	resp := string(out)

	r := regexp.MustCompile(`version: ([0-9]{1}.[0-9]{1}.[0-9]{1})`)
	matches := r.FindAllStringSubmatch(resp, -1)
	log.Printf("Version: %s", matches[1])
}
