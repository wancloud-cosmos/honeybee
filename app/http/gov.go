package http

import (
	"fmt"

	"strconv"
	"time"

	"github.com/astaxie/beego"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	voteUri           = "/gov/proposals/%d/votes"
	queryVoteUri      = "/gov/proposals/%d/votes/%s"
	queryProposalUri  = "/gov/proposals/%d"
	queryProposalsUri = "/gov/proposals"
)

type baseReq struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	ChainID       string `json:"chain_id"`
	AccountNumber int64  `json:"account_number"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
}

type voteReq struct {
	BaseReq baseReq        `json:"base_req"`
	Voter   sdk.AccAddress `json:"voter"`  //  address of the voter
	Option  VoteOption     `json:"option"` //  option from OptionSet chosen by the voter
}

func VoteUrl(id uint64) string {
	return fmt.Sprintf(ClientRESTServerAddress+voteUri, id)
}

func Vote(id uint64, voter string, option VoteOption) error {
	var br baseReq
	br.ChainID = flagChainID
	br.Name = flagName
	br.Password = flagPassword

	aInfo, err := AccountInfo(voter)
	if nil != err {
		beego.Error(err)
		return err
	}
	br.AccountNumber = aInfo.AccountNumber
	br.Sequence = aInfo.Sequence
	br.Gas = 200000 //TODO configable

	var vr voteReq
	addr, err := sdk.AccAddressFromBech32(voter)
	if nil != err {
		beego.Error(err)
		return err
	}
	vr.Voter = addr

	vr.Option = option
	vr.BaseReq = br

	err = POST(VoteUrl(id), vr, nil)
	if nil != err {
		beego.Error(err)
		return err
	}

	return nil
}

type queryVoteResp struct {
	Voter      string `json:"voter"`
	ProposalID string `json:"proposal_id"`
	Ooption    string `json:"option"`
}

func QueryVoteUrl(id uint64, voter string) string {
	return fmt.Sprintf(ClientRESTServerAddress+queryVoteUri, id, voter)
}

func QueryVote(id uint64, voter string) (*queryVoteResp, error) {
	var resp queryVoteResp
	err := GET(QueryVoteUrl(id, voter), &resp)
	if nil != err {
		beego.Error(err)
		return nil, err
	}
	return &resp, nil
}

type queryProposalResp struct {
	Value *queryProposalRespValue `json:"value"`
}

type queryProposalRespValue struct {
	// VotingStartBlock string `json:"voting_start_block"`
	// SubmitBlock      string `json:"submit_block"`
	ProposalID      string    `json:"proposal_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ProposalStatus  string    `json:"proposal_status"`
	VotingStartTime time.Time `json:"voting_start_time"`
}

type Proposal struct {
	// VotingStartBlock int64     `json:"voting_start_block"`
	// SubmitBlock     int64     `json:"submit_block"`
	ProposalID      uint64    `json:"proposal_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ProposalStatus  string    `json:"proposal_status"`
	VotingStartTime time.Time `json:"voting_start_time"`
}

func (p *Proposal) IsPassedStatus() bool {
	return "Passsed" == p.ProposalStatus
}

func (p *Proposal) IsRejectedStatus() bool {
	return "Rejected" == p.ProposalStatus
}

func (p *Proposal) IsDepositPeriodStatus() bool {
	return "DepositPeriod" == p.ProposalStatus
}

func (p *Proposal) IsVotingPeriodStatus() bool {
	return "VotingPeriod" == p.ProposalStatus
}

func QueryProposalUrl(id uint64) string {
	return fmt.Sprintf(ClientRESTServerAddress+queryProposalUri, id)
}

func QueryProposal(id uint64) (*Proposal, error) {
	var resp queryProposalResp
	err := GET(QueryProposalUrl(id), &resp)
	if nil != err {
		beego.Error(err)
		return nil, err
	}

	var p Proposal
	// p.SubmitBlock, err = strconv.ParseInt(resp.Value.SubmitBlock, 10, 64)
	// if nil != err {
	// 	beego.Error(err)
	// 	return nil, err
	// }

	// p.VotingStartBlock, err = strconv.ParseInt(resp.Value.VotingStartBlock, 10, 64)
	// if nil != err {
	// 	beego.Error(err)
	// 	return nil, err
	// }

	p.Description = resp.Value.Description
	p.Title = resp.Value.Title
	p.ProposalID = id
	p.ProposalStatus = resp.Value.ProposalStatus
	p.VotingStartTime = resp.Value.VotingStartTime

	return &p, nil
}

func QueryProposalsUrl(status, depositer, voter string) string {
	return fmt.Sprintf(ClientRESTServerAddress+queryProposalsUri+"?status=%s&depositer=%s&voter=%s", status, depositer, voter)
}

func QueryProposals(status, depositer, voter string) ([]*Proposal, error) {
	var proposals = []queryProposalResp{}
	err := GET(QueryProposalsUrl(status, depositer, voter), &proposals)
	if nil != err {
		beego.Error(err)
		return nil, err
	}

	ps := make([]*Proposal, 0, len(proposals))
	for _, resp := range proposals {
		var p Proposal
		// p.SubmitBlock, err = strconv.ParseInt(resp.Value.SubmitBlock, 10, 64)
		// if nil != err {
		// 	beego.Error(err)
		// 	return nil, err
		// }

		// p.VotingStartBlock, err = strconv.ParseInt(resp.Value.VotingStartBlock, 10, 64)
		// if nil != err {
		// 	beego.Error(err)
		// 	return nil, err
		// }

		p.ProposalID, err = strconv.ParseUint(resp.Value.ProposalID, 10, 64)
		if nil != err {
			beego.Error(err)
			return nil, err
		}

		p.Description = resp.Value.Description
		p.Title = resp.Value.Title
		p.ProposalStatus = resp.Value.ProposalStatus
		p.VotingStartTime = resp.Value.VotingStartTime

		ps = append(ps, &p)
	}

	return ps, nil
}
