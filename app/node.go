package app

import (
	"fmt"
	"strings"
	"validator-monitor/app/config"
	"validator-monitor/utils"

	"github.com/astaxie/beego"
	// "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

var (
	monitorNode *node
)

type node struct {
	Addr string
	cli  *client.HTTP
}

func init() {
	monitorNode = &node{
		Addr: config.NodeAddr,
		cli:  client.NewHTTP(config.NodeAddr, "/websocket"),
	}
}

func (n *node) String() string {
	return fmt.Sprintf("%s", n.Addr)
}

//check if the validator is in the validator set
func (n *node) CheckValidator(addrs []string) error {
	vset, err := n.cli.Validators(nil)
	if nil != err {
		emailBody := fmt.Sprintf("get validator set failed,node:%s,err:%s", n.String(), err.Error())
		beego.Error(emailBody)
		utils.SendMail("get validatorSet failed", emailBody)

		return err
	}

	for _, a := range addrs {
		if !n.IsInVSet(a, vset.Validators) {
			emailBody := fmt.Sprintf("validator:%s is not in vset via node:%s", a, n.String())
			beego.Error(emailBody)
			err = fmt.Errorf(emailBody)
			utils.SendMail("validator is not in vset", emailBody)
			utils.SendSMS(emailBody)
		}
	}

	return err
}

func (n *node) DidMissBlock(addr string) error {
	commit, err := n.cli.Commit(nil)
	if nil != err {
		emailBody := fmt.Sprintf("query commits failed via node(%s),err:%s", n.String(), err.Error())
		beego.Error(emailBody)
		utils.SendMail("query commits failed", emailBody)

		return err
	}

	if !n.IsInLastCommit(addr, commit.Commit.Precommits) {
		err := fmt.Errorf("addr(%s) miss block(%d)", addr, commit.Height)
		beego.Error(err)

		defaultMissBlocks.SetMiss(commit.Height)

		return err
	}

	return nil
}

func (n *node) IsInLastCommit(addr string, commits []*types.Vote) bool {
	for _, v := range commits {
		if nil == v {
			continue
		}

		//just compare fingerprint
		if strings.HasPrefix(addr, v.ValidatorAddress.String()) {
			return true
		}
	}

	return false
}

func (n *node) IsInVSet(addr string, vset []*types.Validator) bool {
	for _, v := range vset {
		if v.Address.String() == addr {
			return true
		}
	}

	return false
}
