package app

import (
	"validator-monitor/app/gov"
	"validator-monitor/app/pubsub"

	"github.com/tendermint/tendermint/libs/pubsub/query"
)

func Start() {
	Watch()

	querySubmitProposal := query.MustParse(`action  = 'submit-proposal'`)
	pubsub.Subscribe(querySubmitProposal, gov.ReadyForVoteHandler)

	queryDeposit := query.MustParse(`action  = 'deposit'`)
	pubsub.Subscribe(queryDeposit, gov.ReadyForVoteHandler)
}
