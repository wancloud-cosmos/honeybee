package app

import (
	"github.com/tendermint/tendermint/libs/pubsub/query"
)

func Start() {
	Watch()

	querySubmitProposal := query.MustParse(`action  = 'submit-proposal'`)
	Subscribe(querySubmitProposal, ReadyForVoteHandler)

	queryDeposit := query.MustParse(`action  = 'deposit'`)
	Subscribe(queryDeposit, ReadyForVoteHandler)
}
