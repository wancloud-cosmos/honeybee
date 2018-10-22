package http

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
)

var (
	accountsUri = "/accounts/%s"
)

type AccountResp struct {
	Value *AccountValueResp `json: "value"`
}

type AccountValueResp struct {
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

type Account struct {
	AccountNumber int64 `json:"account_number"`
	Sequence      int64 `json:"sequence"`
}

func accountsUrl(addr string) string {
	return fmt.Sprintf(ClientRESTServerAddress+accountsUri, addr)
}

func AccountInfo(addr string) (*Account, error) {
	var resp AccountResp
	err := GET(accountsUrl(addr), &resp)
	if nil != err {
		beego.Error(err)
		return nil, err
	}

	var acc Account
	acc.AccountNumber, err = strconv.ParseInt(resp.Value.AccountNumber, 10, 64)
	if nil != err {
		beego.Error(err)
		return nil, err
	}

	acc.Sequence, err = strconv.ParseInt(resp.Value.Sequence, 10, 64)
	if nil != err {
		beego.Error(err)
		return nil, err
	}

	return &acc, nil
}
