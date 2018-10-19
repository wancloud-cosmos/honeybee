package http

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"validator-monitor/utils"

	"github.com/astaxie/beego"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ClientRESTServerAddress string
	GovVoter                string

	flagName     string
	flagPassword string
	flagChainID  string
)

func init() {
	ClientRESTServerAddress = beego.AppConfig.String("client-rest-server:address")
	if "" == ClientRESTServerAddress {
		log.Fatalln("client-rest-server:address invalid")
	}

	GovVoter = beego.AppConfig.String("gov:voter")
	if "" == GovVoter {
		log.Fatalln("gov:voter invalid")
	}

	flagName = beego.AppConfig.String("gov:name")
	if "" == flagName {
		log.Fatalln("gov:name invalid")
	}

	flagPassword = beego.AppConfig.String("gov:password")
	if "" == flagPassword {
		log.Fatalln("gov:password invalid")
	}

	flagChainID = beego.AppConfig.String("gov:chaindid")
	if "" == flagChainID {
		log.Fatalln("gov:chaindid invalid")
	}
}

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

// Type that represents VoteOption as a byte
type VoteOption byte

//nolint
const (
	OptionEmpty      VoteOption = 0x00
	OptionYes        VoteOption = 0x01
	OptionAbstain    VoteOption = 0x02
	OptionNo         VoteOption = 0x03
	OptionNoWithVeto VoteOption = 0x04
)

var (
	voteUri = "/gov/proposals/%d/votes"
)

func VoteUrl(id uint64) string {
	return fmt.Sprintf(ClientRESTServerAddress+voteUri, id)
}

func Vote(id uint64, voter, name, password string, option byte) error {
	var vr voteReq
	vr.Option = OptionYes //TODO

	bin, err := utils.CDC.MarshalJSON(vr)
	if nil != err {
		beego.Error(err)
		return err
	}

	req, err := http.NewRequest("POST", VoteUrl(id), bytes.NewReader(bin))
	if nil != err {
		beego.Error(err)
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		beego.Error(err)
		return err
	}

	if http.StatusOK != resp.StatusCode {
		err = fmt.Errorf("post: %s,expected 200,got:%d", resp.StatusCode)
		beego.Error(err)

		return err
	}

	//TODO notify admin by sending email
	return nil
}
