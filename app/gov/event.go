package gov

import (
	"fmt"
	"validator-monitor/utils"

	"github.com/astaxie/beego"
	"github.com/tendermint/tendermint/libs/events"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	IsGovAutoVote bool
)

func init() {
	IsGovAutoVote = beego.AppConfig.DefaultBool("gov::autovote", false)
}

func ReadyForVoteHandler(query string, data events.EventData) error {
	tags := data.(tmtypes.EventDataTx).TxResult.Result.Tags

	for _, v := range tags {
		beego.Debug(string(v.Key), string(v.Value))

		if "voting-period-start" == string(v.Key) {
			var id int64
			err := utils.CDC.UnmarshalBinaryBare(v.Value, &id)
			if nil != err {
				beego.Error(err)
				return err
			}

			//notify admin user that a new proposal is ready for vote
			subject := fmt.Sprintf("proposal-id:%d ready for vote", id)
			body := subject
			beego.Info(subject)
			utils.SendMail3Times(subject, body)
			utils.SendSMS(subject)

			//auto vote
			if IsGovAutoVote {
				Vote(id)
			}

			return nil
		}
	}

	return nil
}
