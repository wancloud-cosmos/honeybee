package app

import (
	"fmt"
	"time"
	"validator-monitor/app/http"

	"github.com/astaxie/beego"
	"github.com/tendermint/tendermint/rpc/client"
)

var (
	GovVoter string
)

func init() {
	GovVoter = beego.AppConfig.String("gov::voter")
	if "" == GovVoter {
		panic("gov::voter invalid")
	}
}

func LatestBlockHeight() (int64, error) {
	status, err := client.NewHTTP(monitorNodes[0].Addr, "/websocket").Status()
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
				h > p.VotingStartBlock+100 &&
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

	for i := 0; i < 3; i++ {
		err := http.Vote(id, GovVoter, http.OptionNoWithVeto)
		if nil != err {
			beego.Error(err)

			emailTitle += " failed"
			emailBody += " failed,err:" + err.Error()

			time.Sleep(time.Second * 5)
			continue
		} else {
			emailTitle += " success"
			emailBody += " success"

			break
		}
	}

	SendMail3Times(emailTitle, emailBody)

	return nil
}
