package targets

import (
	"bytes"
	"chainflow-vitwit/config"
	client "github.com/influxdata/influxdb1-client/v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type targetRunner struct{}

func NewRunner() *targetRunner {
	return &targetRunner{}
}

func (m targetRunner) Run(function func(ops HTTPOptions, cfg *config.Config, c client.Client), ops HTTPOptions, cfg *config.Config, c client.Client) {
	function(ops, cfg, c)
}

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
				Endpoint: cfg.LCDEndpoint + "bank/balances/" + cfg.OperatorAddress,
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
			Name:          "Deposit Period Proposals",
			HTTPOptions: HTTPOptions{
				Endpoint: "https://api.cosmos.network/gov/proposals",
				Method:   http.MethodGet,
				QueryParams:QueryParams{"status": "deposit_period"},
			},
			Func: GetDepositPeriodProposals,
		},
		{
			ExecutionType: "http",
			Name:          "Voting Period Proposals",
			HTTPOptions: HTTPOptions{
				Endpoint: "https://api.cosmos.network/gov/proposals",
				Method:   http.MethodGet,
				QueryParams:QueryParams{"status": "voting_period"},
			},
			Func: GetVotingPeriodProposals,
		},
		{
			ExecutionType: "http",
			Name:          "Passed Proposals",
			HTTPOptions: HTTPOptions{
				Endpoint: "https://api.cosmos.network/gov/proposals",
				Method:   http.MethodGet,
				QueryParams:QueryParams{"status": "passed"},
			},
			Func: GetPassedProposals,
		},
		{
			ExecutionType: "http",
			Name:          "Rejected Proposals",
			HTTPOptions: HTTPOptions{
				Endpoint: "https://api.cosmos.network/gov/proposals",
				Method:   http.MethodGet,
				QueryParams:QueryParams{"status": "rejected"},
			},
			Func: GetRejectedProposals,
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

func newHttpRequest(ops HTTPOptions) (*http.Request, error) {
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

func HitHTTPTarget(ops HTTPOptions) (*PingResp, error) {
	req, err := newHttpRequest(ops)
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
