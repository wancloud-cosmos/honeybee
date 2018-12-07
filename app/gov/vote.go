package gov

import (
	"fmt"
	"time"
	"validator-monitor/app/config"
	"validator-monitor/app/http"
	"validator-monitor/utils"

	"github.com/astaxie/beego"
	"github.com/tendermint/tendermint/rpc/client"
)

var (
	interval = config.Interval

	GovVoter string

	GovOption http.VoteOption

	GovDelayBlock int64
)

func init() {
	GovVoter = beego.AppConfig.String("gov::voter")
	if "" == GovVoter {
		panic("gov::voter invalid")
	}

	var err error
	GovDelayBlock, err = beego.AppConfig.Int64("gov::delayblock")
	if nil != err {
		panic("gov::delayblock invalid")
	}

	option := beego.AppConfig.String("gov::option")
	if "" == option {
		panic("gov::option invalid")
	}

	GovOption, err = http.VoteOptionFromString(option)
	if nil != err {
		panic("gov::option invalid")
	}
}

func Option(op string) {

}

func LatestBlockHeight() (int64, error) {
	status, err := client.NewHTTP(config.NodeAddr, "/websocket").Status()
	if nil != err {
		beego.Error(err)
		return -1, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

func IsVoted(id int64, voter string) bool {
	_, err := http.QueryVote(id, voter)
	if nil != err {
		beego.Error(err)
		return false
	}

	return true
}

func Vote(id int64) {
	go func() {
		for {
			if IsVoted(id, GovVoter) {
				return
			}

			p, err := http.QueryProposal(id)
			if nil != err {
				beego.Error(err)
				time.Sleep(interval)
				continue
			}

			if p.IsPassedStatus() || p.IsRejectedStatus() {
				return
			}

			h, err := LatestBlockHeight()
			if nil != err {
				beego.Error(err)
				time.Sleep(interval)
				continue
			}

			if p.IsVotingPeriodStatus() &&
				h > p.VotingStartBlock+GovDelayBlock &&
				!IsVoted(id, GovVoter) {

				vote(id)
				return
			}

			time.Sleep(interval)
		}
	}()
}

func vote(id int64) error {
	emailTitle := fmt.Sprintf("vote proposal-id:%d ", id)
	emailBody := fmt.Sprintf("vote proposal-id:%d ", id)

	var err error = nil
	for i := 0; i < 3; i++ {
		err = http.Vote(id, GovVoter, GovOption)
		if nil != err {
			beego.Error(err)
			time.Sleep(time.Second * 5)
			continue
		} else {
			break
		}
	}

	if nil != err {
		emailTitle += " failed"
		emailBody += " failed,err:" + err.Error()
	} else {
		emailTitle += " success"
		emailBody += " success"
	}

	utils.SendMail3Times(emailTitle, emailBody)

	return nil
}
