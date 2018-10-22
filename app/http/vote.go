package http

import (
	"fmt"

	"github.com/astaxie/beego"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	voteUri = "/gov/proposals/%d/votes"
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

func VoteUrl(id int64) string {
	return fmt.Sprintf(ClientRESTServerAddress+voteUri, id)
}

func Vote(id int64, voter string, option VoteOption) error {
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
	br.Sequence = aInfo.Sequence + 1
	br.Gas = 20000 //TODO configable

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
