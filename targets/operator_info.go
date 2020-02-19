package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"strconv"
)

func GetOperatorInfo(ops HTTPOptions, cfg *config.Config, c client.Client) {
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var validatorResp ValidatorResp
	err = json.Unmarshal(resp.Body, &validatorResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	operatorAddress := validatorResp.Validator.Details.OperatorAddress

	address := validatorResp.Validator.Uptime.Address

	fee := validatorResp.Validator.Details.Commission.Rate

	validatorDetails := validatorResp.Validator.Details.Description

	var maxRate float64
	mr, err := strconv.ParseFloat(validatorResp.Validator.Details.Commission.MaxRate, 64)
	if err != nil {
		maxRate = 0
	} else {
		maxRate = mr * 100
	}

	var maxChangeRate float64
	mcr, err := strconv.ParseFloat(validatorResp.Validator.Details.Commission.MaxChangeRate, 64)
	if err != nil {
		log.Printf("error in atoi: %v", err)
		maxChangeRate = 0
	} else {
		maxChangeRate = mcr * 100
	}

	log.Printf("Ooperator Addr: %s \nAddress: %s \nFee: %s \nValidator Details: %v \nMax Rate: %f \nMax Change Rate: %f \n",
		operatorAddress, address, fee, validatorDetails, maxRate, maxChangeRate)
}
