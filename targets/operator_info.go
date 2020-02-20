package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"strconv"
)

func GetOperatorInfo(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}
	var pts []*client.Point

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
	p1, err := createDataPoint("vcf_operator_address", map[string]string{}, map[string]interface{}{"address": operatorAddress})
	if err == nil {
		pts = append(pts, p1)
	}

	address := validatorResp.Validator.Uptime.Address
	p2, err := createDataPoint("vcf_address", map[string]string{}, map[string]interface{}{"address": address})
	if err == nil {
		pts = append(pts, p2)
	}

	fee := validatorResp.Validator.Details.Commission.Rate
	p3, err := createDataPoint("vcf_validator_fee", map[string]string{}, map[string]interface{}{"rate": fee})
	if err == nil {
		pts = append(pts, p3)
	}

	validatorDetails := validatorResp.Validator.Details.Description
	p4, err := createDataPoint("vcf_validator_desc", map[string]string{"tag": "moniker"}, map[string]interface{}{"val": validatorDetails.Moniker})
	p7, err := createDataPoint("vcf_validator_desc", map[string]string{"tag": "website"}, map[string]interface{}{"val": validatorDetails.Website})
	p8, err := createDataPoint("vcf_validator_desc", map[string]string{"tag": "details"}, map[string]interface{}{"val": validatorDetails.Details})
	p9, err := createDataPoint("vcf_validator_desc", map[string]string{"tag": "identity"}, map[string]interface{}{"val": validatorDetails.Identity})
	if err == nil {
		pts = append(pts, p4, p7, p8, p9)
	}

	var maxRate float64
	mr, err := strconv.ParseFloat(validatorResp.Validator.Details.Commission.MaxRate, 64)
	if err != nil {
		maxRate = 0
	} else {
		maxRate = mr * 100
	}
	p5, err := createDataPoint("vcf_validator_max_rate", map[string]string{}, map[string]interface{}{"rate": maxRate})
	if err == nil {
		pts = append(pts, p5)
	}

	var maxChangeRate float64
	mcr, err := strconv.ParseFloat(validatorResp.Validator.Details.Commission.MaxChangeRate, 64)
	if err != nil {
		log.Printf("error in atoi: %v", err)
		maxChangeRate = 0
	} else {
		maxChangeRate = mcr * 100
	}
	p6, err := createDataPoint("vcf_validator_max_change_rate", map[string]string{}, map[string]interface{}{"rate": maxChangeRate})
	if err == nil {
		pts = append(pts, p6)
	}

	bp.AddPoints(pts)
	_ = writeBatchPoints(c, bp)
	log.Printf("Ooperator Addr: %s \nAddress: %s \nFee: %s \nValidator Details: %v \nMax Rate: %f \nMax Change Rate: %f \n",
		operatorAddress, address, fee, validatorDetails, maxRate, maxChangeRate)
}
