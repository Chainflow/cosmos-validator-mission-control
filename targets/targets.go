package targets

import (
	"bytes"
	"chainflow-vitwit/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type targetRunner struct{}

// NewRunner returns targetRunner
func NewRunner() *targetRunner {
	return &targetRunner{}
}

// Run to run the request
func (m targetRunner) Run(function func(ops HTTPOptions, cfg *config.Config, c client.Client), ops HTTPOptions, cfg *config.Config, c client.Client) {
	function(ops, cfg, c)
}

// InitTargets which returns the targets
// can write all the endpoints here
func InitTargets(cfg *config.Config) *Targets {
	return &Targets{List: []Target{
		{
			ExecutionType: "http",
			Name:          "Base URL Endpoint",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.NodeURL,
				Method:   http.MethodGet,
			},
			Func: CheckGaiad,
		},
		{
			ExecutionType: "http",
			Name:          "Net Info URL",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.NodeURL + "net_info?",
				Method:   http.MethodGet,
			},
			Func: GetNetInfo,
		},
		{
			ExecutionType: "cmd",
			Name:          "Gaiacli status cmd",
			Func:          GetGaiaCliStatus,
		},
		{
			ExecutionType: "http",
			Name:          "Operator Information",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "staking/validators/" + cfg.OperatorAddress,
				Method:   http.MethodGet,
			},
			Func: GetOperatorInfo,
		},
		{
			ExecutionType: "http",
			Name:          "Operator Account Information",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "bank/balances/" + cfg.AccountAddress,
				Method:   http.MethodGet,
			},
			Func: GetAccountInfo,
		},
		{
			ExecutionType: "cmd",
			Name:          "Gaiad Version",
			Func:          GaiadVersion,
		},
		{
			ExecutionType: "http",
			Name:          "Proposals",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "gov/proposals",
				Method:   http.MethodGet,
			},
			Func: GetProposals,
		},
		{
			ExecutionType: "http",
			Name:          "Self Delegation",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/staking/delegators/" + cfg.AccountAddress +
					"/delegations/" + cfg.OperatorAddress,
				Method: http.MethodGet,
			},
			Func: GetSelfDelegation,
		},
		{
			ExecutionType: "http",
			Name:          "Current Rewards Amount",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "distribution/validators/" + cfg.OperatorAddress + "/rewards",
				Method:   http.MethodGet,
			},
			Func: GetCurrentRewardsAmount,
		},
		{
			ExecutionType: "http",
			Name:          "Last proposed block and time",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "blocks/latest",
				Method:   http.MethodGet,
			},
			Func: GetLatestProposedBlockAndTime,
		},
		{
			ExecutionType: "cmd",
			Name:          "Latency",
			Func:          GetLatency,
		},
		{
			ExecutionType: "http",
			Name:          "Network Latest Block",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ExternalRPC + "status?",
				Method:   http.MethodGet,
			},
			Func: GetNetworkLatestBlock,
		},
		{
			ExecutionType: "http",
			Name:          "Validator Voting Power",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.NodeURL + "validators",
				Method:   http.MethodGet,
			},
			Func: GetValidatorVotingPower,
		},
		{
			ExecutionType: "http",
			Name:          "Block Time Difference",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.NodeURL + "block",
				Method:   http.MethodGet,
			},
			Func: GetBlockTimeDifference,
		},
		{
			ExecutionType: "http",
			Name:          "Get Current Block Height",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ExternalRPC + "status",
				Method:   http.MethodGet,
			},
			Func: GetMissedBlocks,
		},
	}}
}

func addQueryParameters(req *http.Request, queryParams QueryParams) {
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	req.URL.RawQuery = params.Encode()
}

// newHTTPRequest to make a new http request
func newHTTPRequest(ops HTTPOptions) (*http.Request, error) {
	// make new request
	req, err := http.NewRequest(ops.Method, ops.Endpoint, bytes.NewBuffer(ops.Body))
	if err != nil {
		return nil, err
	}

	// Add any query parameters to the URL.
	if len(ops.QueryParams) != 0 {
		addQueryParameters(req, ops.QueryParams)
	}

	return req, nil
}

func makeResponse(res *http.Response) (*PingResp, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &PingResp{}, err
	}

	response := &PingResp{
		StatusCode: res.StatusCode,
		Body:       body,
	}
	_ = res.Body.Close()
	return response, nil
}

// HitHTTPTarget to hit the target and get response
func HitHTTPTarget(ops HTTPOptions) (*PingResp, error) {
	req, err := newHTTPRequest(ops)
	if err != nil {
		return nil, err
	}

	httpcli := http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := httpcli.Do(req)
	if err != nil {
		return nil, err
	}

	res, err := makeResponse(resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}
