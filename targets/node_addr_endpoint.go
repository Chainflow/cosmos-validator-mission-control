package targets

import "log"

func GetNodeAddrEndpointData(ops HTTPOptions) {
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.StatusCode == 200 {
		log.Println("Gaiad is running...")
		return
	}

	log.Printf("Error response from gaiad: %s", string(resp.Body))
}
