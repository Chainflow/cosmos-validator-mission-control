package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetValidatorVoted to check validator voted for the proposal or not
func GetValidatorVoted(LCDEndpoint string, proposalID string, accountAddress string) string {
	proposalURL := LCDEndpoint + "/cosmos/gov/v1beta1/proposals/" + proposalID + "/votes"
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
	for _, value := range voters.Votes {
		if value.Voter == accountAddress {
			validatorVoted = returnVoteOption(value.Option)
		}
	}
	return validatorVoted
}

func returnVoteOption(option string) string {
	m := map[string]string{
		"VOTE_OPTION_YES":          "Yes",
		"VOTE_OPTION_ABSTAIN":      "Abstain",
		"VOTE_OPTION_NO":           "No",
		"VOTE_OPTION_NO_WITH_VETO": "NoWithVeto",
	}

	value := m[option]
	return value
}

// SendVotingPeriodProposalAlerts which send alerts of voting period proposals
func SendVotingPeriodProposalAlerts(LCDEndpoint string, accountAddress string, cfg *config.Config, c client.Client) error {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	proposalURL := LCDEndpoint + "/cosmos/gov/v1beta1/proposals?status=PROPOSAL_STATUS_VOTING_PERIOD"
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

	for _, proposal := range p.Proposals {
		proposalVotesURL := LCDEndpoint + "/cosmos/gov/v1beta1/proposals/" + proposal.ProposalID + "/votes"
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
		for _, value := range voters.Votes {
			if value.Voter == accountAddress {
				validatorVoted = returnVoteOption(value.Option)
			}
		}

		if validatorVoted == "No" {
			now := time.Now().UTC()
			votingEndTime, _ := time.Parse(time.RFC3339, proposal.VotingEndTime)
			timeDiff := now.Sub(votingEndTime)
			log.Println("timeDiff...", timeDiff.Hours())

			var proposalAlertCount = 1
			count := GetVotesProposalAlertsCount(cfg, c, proposal.ProposalID)
			if count != "" {
				pac, err := strconv.Atoi(count)
				if err != nil {
					log.Printf("Error while converting proposal alerts count : %v", err)
					return err
				}
				proposalAlertCount = pac
			}

			if timeDiff.Hours() == 24 && proposalAlertCount <= 1 {
				_ = SendTelegramAlert(fmt.Sprintf("%s validator has not voted on proposal = %s, No.of hours left to vote is : %f", cfg.ValidatorName, proposal.ProposalID, timeDiff.Hours()), cfg)
				_ = SendEmailAlert(fmt.Sprintf("%s validator has not voted on proposal = %s", cfg.ValidatorName, proposal.ProposalID), cfg)

				proposalAlertCount = proposalAlertCount + 1
				_ = writeToInfluxDb(c, bp, "vcf_votes_proposal_alert_count", map[string]string{}, map[string]interface{}{"count": proposalAlertCount, "proposal_id": proposal.ProposalID})

				log.Println("Sent alert of voting period proposals")
			}

			if timeDiff.Hours() == 12 && proposalAlertCount <= 2 {
				_ = SendTelegramAlert(fmt.Sprintf("%s validator has not voted on proposal = %s, No.of hours left to vote is : %f", cfg.ValidatorName, proposal.ProposalID, timeDiff.Hours()), cfg)
				_ = SendEmailAlert(fmt.Sprintf("%s validator has not voted on proposal = %s", cfg.ValidatorName, proposal.ProposalID), cfg)

				proposalAlertCount = proposalAlertCount + 1
				_ = writeToInfluxDb(c, bp, "vcf_votes_proposal_alert_count", map[string]string{}, map[string]interface{}{"count": proposalAlertCount, "proposal_id": proposal.ProposalID})

				log.Println("Sent alert of voting period proposals")
			}
		}
	}
	return nil
}

// GetValidatorDeposited to check validator deposited for the proposal or not
func GetValidatorDeposited(LCDEndpoint string, proposalID string, accountAddress string) string {
	proposalURL := LCDEndpoint + "/cosmos/gov/v1beta1/proposals/" + proposalID + "/deposits"
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
	for _, value := range depositors.Deposits {
		if value.Depositor == accountAddress && len(value.Amount) != 0 {
			validateDeposit = "yes"
			break
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

	for _, proposal := range p.Proposals {
		validatorVoted := GetValidatorVoted(cfg.LCDEndpoint, proposal.ProposalID, cfg.AccountAddress)
		validatorDeposited := GetValidatorDeposited(cfg.LCDEndpoint, proposal.ProposalID, cfg.AccountAddress)
		err = SendVotingPeriodProposalAlerts(cfg.LCDEndpoint, cfg.AccountAddress, cfg, c)
		if err != nil {
			log.Printf("Error while sending voting period alert: %v", err)
		}

		tag := map[string]string{"id": proposal.ProposalID}
		fields := map[string]interface{}{
			"proposal_id":               proposal.ProposalID,
			"content_type":              proposal.Content.Type,
			"content_value_title":       proposal.Content.Title,
			"content_value_description": proposal.Content.Description,
			"proposal_status":           proposal.Status,
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
		q := client.NewQuery(fmt.Sprintf("SELECT * FROM vcf_proposals WHERE proposal_id = '%s'", proposal.ProposalID), cfg.InfluxDB.Database, "")
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
				log.Printf("New Proposal Came: %s", proposal.ProposalID)
				_ = writeToInfluxDb(c, bp, "vcf_proposals", tag, fields)

				if proposal.Status == "PROPOSAL_STATUS_REJECTED" || proposal.Status == "PROPOSAL_STATUS_PASSED" { // proposal status
					_ = SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been %s", proposal.ProposalID, proposal.Status), cfg)
					_ = SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been = %s", proposal.ProposalID, proposal.Status), cfg)
				} else if proposal.Status == "PROPOSAL_STATUS_VOTING_PERIOD" {
					_ = SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ProposalID, proposal.Status), cfg)
					_ = SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ProposalID, proposal.Status), cfg)
				} else {
					_ = SendTelegramAlert(fmt.Sprintf("A new proposal "+proposal.Content.Type+" has been added to "+proposal.Status+" with proposal id = %s", proposal.ProposalID), cfg)
					_ = SendEmailAlert(fmt.Sprintf("A new proposal "+proposal.Content.Type+" has been added to "+proposal.Status+" with proposal id = %s", proposal.ProposalID), cfg)
				}
			} else {
				q := client.NewQuery(fmt.Sprintf("DELETE FROM vcf_proposals WHERE id = '%s'", proposal.ProposalID), cfg.InfluxDB.Database, "")
				if response, err := c.Query(q); err == nil && response.Error() == nil {
					log.Printf("Delete proposal %s from vcf_proposals", proposal.ProposalID)
				} else {
					log.Printf("Failed to delete proposal %s from vcf_proposals", proposal.ProposalID)
					log.Printf("Reason for proposal deletion failure %v", response)
				}
				log.Printf("Writing the proposal: %s", proposal.ProposalID)
				_ = writeToInfluxDb(c, bp, "vcf_proposals", tag, fields)
				if proposal.Status != proposalStatus {
					if proposal.Status == "PROPOSAL_STATUS_REJECTED" || proposal.Status == "PROPOSAL_STATUS_PASSED" {
						_ = SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been %s", proposal.ProposalID, proposal.Status), cfg)
						_ = SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been = %s", proposal.ProposalID, proposal.Status), cfg)
					} else {
						_ = SendTelegramAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ProposalID, proposal.Status), cfg)
						_ = SendEmailAlert(fmt.Sprintf("Proposal "+proposal.Content.Type+" with proposal id = %s has been moved to %s", proposal.ProposalID, proposal.Status), cfg)
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

					for _, proposal := range p.Proposals {
						if proposal.ProposalID == ID {
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

// GetVotesProposalAlertsCount returns the count of voting period alerts
func GetVotesProposalAlertsCount(cfg *config.Config, c client.Client, proposalID string) string {
	var count string
	q := client.NewQuery(fmt.Sprintf("SELECT last(count) FROM vcf_votes_proposal_alert_count WHERE proposal_id = '%s'", proposalID), cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						pc := r.Series[0].Values[0][idx]
						count = fmt.Sprintf("%v", pc)
						break
					}
				}
			}
		}
	}

	return count
}
