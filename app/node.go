package app

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

var (
	monitorNodes []*node
)

func init() {
	nodeAddresses := strings.Split(beego.AppConfig.String("node::address"), ",")
	for _, v := range nodeAddresses {
		n := node{Addr: v}
		monitorNodes = append(monitorNodes, &n)
	}

	beego.Debug("config node::address ", nodeAddresses)
}

type node struct {
	Addr string
}

func (n *node) String() string {
	return fmt.Sprintf("%s", n.Addr)
}

//check if the validator in the validator set
func (n *node) CheckValidator(addrs []string) {
	vset, err := client.NewHTTP(n.Addr, "/websocket").Validators(nil)
	if nil != err {
		emailBody := fmt.Sprintf("get validator set failed,node:%s,err:%s", n.String(), err.Error())
		beego.Error(emailBody)
		SendMail(emailTos, "get validator set failed", emailBody)

		return //TODO try 3 times?
	}

	for _, a := range addrs {
		if !n.IsInVSet(a, vset.Validators) {
			emailBody := fmt.Sprintf("validator:%s is not in vset via node:%s", a, n.String())
			beego.Error(emailBody)
			SendMail(emailTos, "validator is not in vset", emailBody)
		}
	}
}

func (n *node) IsInVSet(addr string, vset []*types.Validator) bool {
	for _, v := range vset {
		if v.Address.String() == addr {
			return true
		}
	}

	return false
}
