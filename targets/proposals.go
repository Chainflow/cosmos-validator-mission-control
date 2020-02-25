package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

func GetDepositPeriodProposals(ops HTTPOptions, cfg *config.Config, c client.Client) {
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
		tag := map[string]string{"id": proposal.Id}
		fields := map[string]interface{}{
			"content.type":       proposal.Content.Type,
			"content.value":      proposal.Content.Value,
			"proposal_status":    proposal.ProposalStatus,
			"final_tally_result": proposal.FinalTallyResult,
			"submit_time":        proposal.SubmitTime,
			"deposit_end_time":   proposal.DepositEndTime,
			"total_deposit":      proposal.TotalDeposit,
			"voting_start_time":  proposal.VotingStartTime,
			"voting_end_time":    proposal.VotingEndTime,
		}
		newProposal := false
		q := client.NewQuery(fmt.Sprintf("SELECT count(*) as count FROM vcf_deposit_period_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
		if response, err := c.Query(q); err == nil && response.Error() == nil {
			for _, r := range response.Results {
				if len(r.Series) == 0 {
					newProposal = true
					break
				}
			}

			if newProposal {
				log.Printf("New Proposal Came In Deposit Period: %s", proposal.Id)
				_ = writeToInfluxDb(c, bp, "vcf_deposit_period_proposals", tag, fields)
				_ = SendTelegramAlert(fmt.Sprintf("A new proposal has been added to deposit period with proposal id = %s", proposal.Id), cfg)
				_ = SendEmailAlert(fmt.Sprintf("A new proposal has been added to deposit period with proposal id = %s", proposal.Id), cfg)
			}
		}
	}
}

func GetVotingPeriodProposals(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var p VotingPeriodProposal
	err = json.Unmarshal(resp.Body, &p)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for _, proposal := range p.Result {
		tag := map[string]string{"id": proposal.Id}
		fields := map[string]interface{}{
			"content.type":       proposal.Content.Type,
			"content.value":      proposal.Content.Value,
			"proposal_status":    proposal.ProposalStatus,
			"final_tally_result": proposal.FinalTallyResult,
			"submit_time":        proposal.SubmitTime,
			"deposit_end_time":   proposal.DepositEndTime,
			"total_deposit":      proposal.TotalDeposit,
			"voting_start_time":  proposal.VotingStartTime,
			"voting_end_time":    proposal.VotingEndTime,
		}
		newProposal := false
		q := client.NewQuery(fmt.Sprintf("SELECT count(*) as count FROM vcf_voting_period_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
		if response, err := c.Query(q); err == nil && response.Error() == nil {
			for _, r := range response.Results {
				if len(r.Series) == 0 {
					newProposal = true
					break
				}
			}

			if newProposal {
				log.Printf("New Proposal Came In Voting Period: %s", proposal.Id)
				_ = writeToInfluxDb(c, bp, "vcf_voting_period_proposals", tag, fields)
				_ = SendTelegramAlert(fmt.Sprintf("A new proposal has been added to voting period with proposal id = %s", proposal.Id), cfg)
				_ = SendEmailAlert(fmt.Sprintf("A new proposal has been added to voting period with proposal id = %s", proposal.Id), cfg)

				q := client.NewQuery(fmt.Sprintf("DELETEE  FROM vcf_deposit_period_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
				if response, err := c.Query(q); err == nil && response.Error() == nil {
					log.Printf("Delete proposal %s from vcf_deposit_period_proposals", proposal.Id)
				} else {
					log.Printf("Failed to delete proposal %s from vcf_deposit_period_proposals", proposal.Id)
				}
			}
		}
	}
}

func GetPassedProposals(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var p PassedProposal
	err = json.Unmarshal(resp.Body, &p)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for _, proposal := range p.Result {
		tag := map[string]string{"id": proposal.Id}
		fields := map[string]interface{}{
			"content.type":       proposal.Content.Type,
			"content.value":      proposal.Content.Value,
			"proposal_status":    proposal.ProposalStatus,
			"final_tally_result": proposal.FinalTallyResult,
			"submit_time":        proposal.SubmitTime,
			"deposit_end_time":   proposal.DepositEndTime,
			"total_deposit":      proposal.TotalDeposit,
			"voting_start_time":  proposal.VotingStartTime,
			"voting_end_time":    proposal.VotingEndTime,
		}
		newProposal := false
		q := client.NewQuery(fmt.Sprintf("SELECT count(*) as count FROM vcf_passed_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
		if response, err := c.Query(q); err == nil && response.Error() == nil {
			for _, r := range response.Results {
				if len(r.Series) == 0 {
					newProposal = true
					break
				}
			}

			if newProposal {
				log.Printf("New Proposal Passed with Proposal ID: %s", proposal.Id)
				_ = writeToInfluxDb(c, bp, "vcf_passed_proposals", tag, fields)
				_ = SendTelegramAlert(fmt.Sprintf("A new proposal has passed with proposal id = %s", proposal.Id), cfg)
				_ = SendEmailAlert(fmt.Sprintf("A new proposal has passeed with proposal id = %s", proposal.Id), cfg)

				q := client.NewQuery(fmt.Sprintf("DELETEE FROM vcf_voting_period_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
				if response, err := c.Query(q); err == nil && response.Error() == nil {
					log.Printf("Delete proposal %s from vcf_voting_period_proposals", proposal.Id)
				} else {
					log.Printf("Failed to delete proposal %s from vcf_voting_period_proposals", proposal.Id)
				}
			}
		}
	}
}

func GetRejectedProposals(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var p RejectedProposal
	err = json.Unmarshal(resp.Body, &p)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for _, proposal := range p.Result {
		tag := map[string]string{"id": proposal.Id}
		fields := map[string]interface{}{
			"content.type":       proposal.Content.Type,
			"content.value":      proposal.Content.Value,
			"proposal_status":    proposal.ProposalStatus,
			"final_tally_result": proposal.FinalTallyResult,
			"submit_time":        proposal.SubmitTime,
			"deposit_end_time":   proposal.DepositEndTime,
			"total_deposit":      proposal.TotalDeposit,
			"voting_start_time":  proposal.VotingStartTime,
			"voting_end_time":    proposal.VotingEndTime,
		}
		newProposal := false
		q := client.NewQuery(fmt.Sprintf("SELECT count(*) as count FROM vcf_rejected_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
		if response, err := c.Query(q); err == nil && response.Error() == nil {
			for _, r := range response.Results {
				if len(r.Series) == 0 {
					newProposal = true
					break
				}
			}

			if newProposal {
				log.Printf("Proposal Rejected with Proposal ID: %s", proposal.Id)
				_ = writeToInfluxDb(c, bp, "vcf_rejected_proposals", tag, fields)
				_ = SendTelegramAlert(fmt.Sprintf("A new proposal has been rejected with proposal id = %s", proposal.Id), cfg)
				_ = SendEmailAlert(fmt.Sprintf("A new proposal has been rejected with proposal id = %s", proposal.Id), cfg)

				q := client.NewQuery(fmt.Sprintf("DELETEE FROM vcf_voting_period_proposals WHERE id = '%s'", proposal.Id), cfg.InfluxDB.Database, "")
				if response, err := c.Query(q); err == nil && response.Error() == nil {
					log.Printf("Delete proposal %s from vcf_voting_period_proposals", proposal.Id)
				} else {
					log.Printf("Failed to delete proposal %s from vcf_voting_period_proposals", proposal.Id)
				}
			}
		}
	}
}
