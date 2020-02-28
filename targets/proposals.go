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

		// Get proposal voters
		proposalURL := "http://http://134.209.142.233:1317/" + "gov/proposals/" + proposal.Id + "/votes"
		res, err := http.Get(proposalURL)
		if err != nil {
			log.Printf("Error: %v", err)
			return
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
			if value.Voter == cfg.AccountAddress {
				validatorVoted = "yes"
			}
		}

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
				log.Printf("New Proposal Came In Deposit Period: %s", proposal.Id)
				_ = writeToInfluxDb(c, bp, "vcf_proposals", tag, fields)
				_ = SendTelegramAlert(fmt.Sprintf("A new proposal has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.Id), cfg)
				_ = SendEmailAlert(fmt.Sprintf("A new proposal has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.Id), cfg)
			} else {
				q := client.NewQuery(fmt.Sprintf("UPDATE vcf_proposals SET proposal_status=%s, validator_voted=%s,voting_start_time=%s,voting_end_time=%s WHERE id = '%s'", proposal.ProposalStatus, validatorVoted, proposal.VotingStartTime, proposal.VotingEndTime, proposal.Id), cfg.InfluxDB.Database, "")
				_, err := c.Query(q)
				if err != nil {
					log.Print("Error while updating proposal ", err)
					return
				}
				_ = SendTelegramAlert(fmt.Sprintf("A new proposal has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.Id), cfg)
				_ = SendEmailAlert(fmt.Sprintf("A new proposal has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.Id), cfg)
			}
		}
	}
}
