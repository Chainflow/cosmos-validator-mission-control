package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	client "github.com/influxdata/influxdb1-client/v2"
)

// Check validator voted for the proposal or not
func GetValidatorVoted(LCDEndpoint string, proposalID string, accountAddress string) string {

	proposalURL := LCDEndpoint + "gov/proposals/" + proposalID + "/votes"
	res, err := http.Get(proposalURL)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	var voters ProposalVoters
	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error while reading resp body ", err)
		}
		_ = json.Unmarshal(body, &voters)
	}

	validatorVoted := "no"
	for _, value := range voters.Result {
		if value.Voter == accountAddress {
			validatorVoted = "yes"
		}
	}
	return validatorVoted
}

func GetProposals(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var p DepositPeriodProposal
	err = json.Unmarshal(resp.Body, &p)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for _, proposal := range p.Result {

		validatorVoted := GetValidatorVoted(cfg.LCDEndpoint, proposal.Id, cfg.AccountAddress)

		tag := map[string]string{"id": proposal.Id}
		fields := map[string]interface{}{
			"id":                        proposal.Id,
			"content_type":              proposal.Content.Type,
			"content_value_title":       proposal.Content.Value.Title,
			"content_value_description": proposal.Content.Value.Description,
			"proposal_status":           proposal.ProposalStatus,
			"final_tally_result":        proposal.FinalTallyResult,
			"submit_time":               proposal.SubmitTime,
			"deposit_end_time":          proposal.DepositEndTime,
			"total_deposit":             proposal.TotalDeposit,
			"voting_start_time":         proposal.VotingStartTime,
			"voting_end_time":           proposal.VotingEndTime,
			"validator_voted":           validatorVoted,
		}
		newProposal := false
		q := client.NewQuery(fmt.Sprintf("SELECT * FROM vcf_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
		if response, err := c.Query(q); err == nil && response.Error() == nil {
			for _, r := range response.Results {
				if len(r.Series) == 0 {
					newProposal = true
					break
				}
			}

			if newProposal {
				log.Printf("New Proposal Came: %s", proposal.Id)
				_ = writeToInfluxDb(c, bp, "vcf_proposals", tag, fields)
				_ = SendTelegramAlert(fmt.Sprintf("A new proposal has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.Id), cfg)
				_ = SendEmailAlert(fmt.Sprintf("A new proposal has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.Id), cfg)
			} else {
				q := client.NewQuery(fmt.Sprintf("DELETE  FROM vcf_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
				if response, err := c.Query(q); err == nil && response.Error() == nil {
					log.Printf("Delete proposal %s from vcf_proposals", proposal.Id)
				} else {
					log.Printf("Failed to delete proposal %s from vcf_proposals", proposal.Id)
				}
				log.Printf("Updating the proposal: %s", proposal.Id)
				_ = writeToInfluxDb(c, bp, "vcf_proposals", tag, fields)
				_ = SendTelegramAlert(fmt.Sprintf("A new proposal has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.Id), cfg)
				_ = SendEmailAlert(fmt.Sprintf("A new proposal has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.Id), cfg)
			}
		}
	}
}
