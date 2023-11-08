package server

type (
	// QueryParams to map the query params of an url
	QueryParams map[string]string

	// HTTPOptions of a target
	HTTPOptions struct {
		Endpoint    string
		QueryParams QueryParams
		Body        []byte
		Method      string
	}

	// PingResp struct
	PingResp struct {
		StatusCode int
		Body       []byte
	}

	// NetworkLatestBlock stores latest block height info
	NetworkLatestBlock struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}

	// CurrentBlock is a struct which holds latest current block details
	CurrentBlock struct {
		Result struct {
			Block struct {
				LastCommit struct {
					Signatures []struct {
						ValidatorAddress string `json:"validator_address"`
						Signature        string `json:"signature"`
					} `json:"signatures"`
				} `json:"last_commit"`
			} `json:"block"`
		} `json:"result"`
	}

	// Validator is a struct which holds validator staking info
	Validator struct {
		Validator struct {
			Jailed bool   `json:"jailed"`
			Status string `json:"status"`
		} `json:"validator"`
	}

	// ValidatorRpcStatus is a struct which holds the status response
	ValidatorRpcStatus struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			SyncInfo struct {
				LatestBlockHash   string `json:"latest_block_hash"`
				LatestAppHash     string `json:"latest_app_hash"`
				LatestBlockHeight string `json:"latest_block_height"`
				LatestBlockTime   string `json:"latest_block_time"`
				CatchingUp        bool   `json:"catching_up"`
			} `json:"sync_info"`
			ValidatorInfo struct {
				VotingPower string `json:"voting_power"`
			} `json:"validator_info"`
		} `json:"result"`
	}

	// ProposalVoters struct holds the parameters of proposal voters
	ProposalVoters struct {
		Votes []struct {
			ProposalID string `json:"proposal_id"`
			Voter      string `json:"voter"`
			Option     string `json:"option"`
		} `json:"votes"`
		Pagination struct {
			NextKey interface{} `json:"next_key"`
			Total   string      `json:"total"`
		} `json:"pagination"`
	}

	// Proposals struct holds result of array of proposals
	Proposals struct {
		Proposals []ProposalResult `json:"proposals"`
	}

	// ProposalResult struct holds the parameters of proposal result
	ProposalResult struct {
		Content          ProposalResultContent `json:"content"`
		ProposalID       string                `json:"proposal_id"`
		Status           string                `json:"status"`
		FinalTallyResult interface{}           `json:"final_tally_result"`
		SubmitTime       string                `json:"submit_time"`
		DepositEndTime   string                `json:"deposit_end_time"`
		TotalDeposit     []interface{}         `json:"total_deposit"`
		VotingStartTime  string                `json:"voting_start_time"`
		VotingEndTime    string                `json:"voting_end_time"`
	}

	// ProposalResultContent struct holds the parameters of a proposal content result
	ProposalResultContent struct {
		Type        string `json:"@type"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	// Depositors struct which holds the parameters of deposits
	Depositors struct {
		Deposits []struct {
			ProposalID string `json:"proposal_id"`
			Depositor  string `json:"depositor"`
			Amount     []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
		} `json:"deposits"`
		Pagination struct {
			NextKey interface{} `json:"next_key"`
			Total   string      `json:"total"`
		} `json:"pagination"`
	}
)
