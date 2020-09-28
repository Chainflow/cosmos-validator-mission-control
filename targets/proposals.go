package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetValidatorVoted to check validator voted for the proposal or not
func GetValidatorVoted(LCDEndpoint string, proposalID string, accountAddress string) string {
	proposalURL := LCDEndpoint + "/gov/proposals/" + proposalID + "/votes"
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

	validatorVoted := "not voted"
	for _, value := range voters.Result {
		if value.Voter == accountAddress {
			validatorVoted = value.Option
		}
	}
	return validatorVoted
}

// SendVotingPeriodProposalAlerts which send alerts of voting period proposals
func SendVotingPeriodProposalAlerts(LCDEndpoint string, accountAddress string, cfg *config.Config) error {
	proposalURL := LCDEndpoint + "/gov/proposals?status=voting_period"
	res, err := http.Get(proposalURL)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	var p Proposals
	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error while reading resp body ", err)
			return err
		}
		_ = json.Unmarshal(body, &p)
	}

	for _, proposal := range p.Result {
		proposalVotesURL := LCDEndpoint + "/gov/proposals/" + proposal.ID + "/votes"
		res, err := http.Get(proposalVotesURL)
		if err != nil {
			log.Printf("Error: %v", err)
			return err
		}

		var voters ProposalVoters
		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error while reading resp body ", err)
				return err
			}
			_ = json.Unmarshal(body, &voters)
		}

		var validatorVoted string
		for _, value := range voters.Result {
			if value.Voter == accountAddress {
				validatorVoted = value.Option
			}
		}

		if validatorVoted == "No" {
			now := time.Now().UTC()
			votingEndTime, _ := time.Parse(time.RFC3339, proposal.VotingEndTime)
			timeDiff := now.Sub(votingEndTime)
			log.Println("timeDiff...", timeDiff.Hours())

			if timeDiff.Hours() <= 24 {
				_ = SendTelegramAlert(fmt.Sprintf("%s validator has not voted on proposal = %s", cfg.ValidatorName, proposal.ID), cfg)
				_ = SendEmailAlert(fmt.Sprintf("%s validator has not voted on proposal = %s", cfg.ValidatorName, proposal.ID), cfg)

				log.Println("Sent alert of voting period proposals")
			}
		}
	}
	return nil
}

// GetValidatorDeposited to check validator deposited for the proposal or not
func GetValidatorDeposited(LCDEndpoint string, proposalID string, accountAddress string) string {

	proposalURL := LCDEndpoint + "/gov/proposals/" + proposalID + "/deposits"
	res, err := http.Get(proposalURL)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	var depositors Depositors
	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error while reading resp body ", err)
		}
		_ = json.Unmarshal(body, &depositors)
	}

	validateDeposit := "no"
	for _, value := range depositors.Result {
		if value.Depositor == accountAddress && len(value.Amount) != 0 {
			validateDeposit = "yes"
		}
	}
	return validateDeposit
}

// GetProposals to get all the proposals and send alerts accordingly
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

	var p Proposals
	err = json.Unmarshal(resp.Body, &p)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for _, proposal := range p.Result {
		validatorVoted := GetValidatorVoted(cfg.LCDEndpoint, proposal.ID, cfg.AccountAddress)
		validatorDeposited := GetValidatorDeposited(cfg.LCDEndpoint, proposal.ID, cfg.AccountAddress)
		err = SendVotingPeriodProposalAlerts(cfg.LCDEndpoint, cfg.AccountAddress, cfg)
		if err != nil {
			log.Printf("Error while sending voting period alert: %v", err)
		}

		tag := map[string]string{"id": proposal.ID}
		fields := map[string]interface{}{
			"proposal_id":               proposal.ID,
			"content_type":              proposal.Content.Type,
			"content_value_title":       proposal.Content.Value.Title,
			"content_value_description": proposal.Content.Value.Description,
			"proposal_status":           proposal.ProposalStatus,
			"final_tally_result":        proposal.FinalTallyResult,
			"submit_time":               GetUserDateFormat(proposal.SubmitTime),
			"deposit_end_time":          GetUserDateFormat(proposal.DepositEndTime),
			"total_deposit":             proposal.TotalDeposit,
			"voting_start_time":         GetUserDateFormat(proposal.VotingStartTime),
			"voting_end_time":           GetUserDateFormat(proposal.VotingEndTime),
			"validator_voted":           validatorVoted,
			"validator_deposited":       validatorDeposited,
		}
		newProposal := false
		proposalStatus := ""
		q := client.NewQuery(fmt.Sprintf("SELECT * FROM vcf_proposals WHERE proposal_id = '%s'", proposal.ID), cfg.InfluxDB.Database, "")
		if response, err := c.Query(q); err == nil && response.Error() == nil {
			for _, r := range response.Results {
				if len(r.Series) == 0 {
					newProposal = true
					break
				} else {
					for idx, col := range r.Series[0].Columns {
						if col == "proposal_status" {
							v := r.Series[0].Values[0][idx]
							proposalStatus = fmt.Sprintf("%v", v)
						}
					}
				}
			}

			if newProposal {
				log.Printf("New Proposal Came: %s", proposal.ID)
				_ = writeToInfluxDb(c, bp, "vcf_proposals", tag, fields)

				if proposal.ProposalStatus == "Rejected" || proposal.ProposalStatus == "Passed" {
					_ = SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been %s", proposal.ID, proposal.ProposalStatus), cfg)
					_ = SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been = %s", proposal.ID, proposal.ProposalStatus), cfg)
				} else if proposal.ProposalStatus == "VotingPeriod" {
					_ = SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ID, proposal.ProposalStatus), cfg)
					_ = SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ID, proposal.ProposalStatus), cfg)
				} else {
					_ = SendTelegramAlert(fmt.Sprintf("A new proposal "+proposal.Content.Type+" has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.ID), cfg)
					_ = SendEmailAlert(fmt.Sprintf("A new proposal "+proposal.Content.Type+" has been added to "+proposal.ProposalStatus+" with proposal id = %s", proposal.ID), cfg)
				}
			} else {
				q := client.NewQuery(fmt.Sprintf("DELETE FROM vcf_proposals WHERE id = '%s'", proposal.ID), cfg.InfluxDB.Database, "")
				if response, err := c.Query(q); err == nil && response.Error() == nil {
					log.Printf("Delete proposal %s from vcf_proposals", proposal.ID)
				} else {
					log.Printf("Failed to delete proposal %s from vcf_proposals", proposal.ID)
					log.Printf("Reason for proposal deletion failure %v", response)
				}
				log.Printf("Writing the proposal: %s", proposal.ID)
				_ = writeToInfluxDb(c, bp, "vcf_proposals", tag, fields)
				if proposal.ProposalStatus != proposalStatus {
					if proposal.ProposalStatus == "Rejected" || proposal.ProposalStatus == "Passed" {
						_ = SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been %s", proposal.ID, proposal.ProposalStatus), cfg)
						_ = SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been = %s", proposal.ID, proposal.ProposalStatus), cfg)
					} else {
						_ = SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ID, proposal.ProposalStatus), cfg)
						_ = SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ID, proposal.ProposalStatus), cfg)
					}
				}
			}
		}
	}

	// Calling fucntion to delete deposit proposals
	// which are ended
	err = DeleteDepoitEndProposals(cfg, c, p)
	if err != nil {
		log.Printf("Error while deleting proposals")
	}
}

// DeleteDepoitEndProposals to delete proposals from db
//which are not present in lcd resposne
func DeleteDepoitEndProposals(cfg *config.Config, c client.Client, p Proposals) error {
	var ID string
	found := false
	q := client.NewQuery("SELECT * FROM vcf_proposals where proposal_status='DepositPeriod'", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx := range r.Series[0].Values {
					proposalID := r.Series[0].Values[idx][7]
					ID = fmt.Sprintf("%v", proposalID)

					for _, proposal := range p.Result {
						if proposal.ID == ID {
							found = true
							break
						} else {
							found = false
						}
					}
					if !found {
						q := client.NewQuery(fmt.Sprintf("DELETE FROM vcf_proposals WHERE id = '%s'", ID), cfg.InfluxDB.Database, "")
						if response, err := c.Query(q); err == nil && response.Error() == nil {
							log.Printf("Delete proposal %s from vcf_proposals", ID)
							return err
						}
						log.Printf("Failed to delete proposal %s from vcf_proposals", ID)
						log.Printf("Reason for proposal deletion failure %v", response)
					}
				}
			}
		}
	}
	return nil
}

// GetUserDateFormat to which returns date in a user friendly
func GetUserDateFormat(timeToConvert string) string {
	time, err := time.Parse(time.RFC3339, timeToConvert)
	if err != nil {
		fmt.Println("Error while converting date ", err)
	}
	date := time.Format("Mon Jan _2 15:04:05 2006")
	fmt.Println("Converted time into date format : ", date)
	return date
}
