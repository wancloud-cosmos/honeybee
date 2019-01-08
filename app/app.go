package app

import (
	"fmt"
	"validator-monitor/app/gov"
	"validator-monitor/app/pubsub"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/tendermint/tendermint/libs/pubsub/query"
)

func Start() {
	Watch()

	querySubmitProposal := query.MustParse(fmt.Sprintf(`action  = '%s'`, gov.TagActionSubmitProposal))
	pubsub.Subscribe(querySubmitProposal, gov.ReadyForVoteHandler)

	queryDeposit := query.MustParse(fmt.Sprintf(`action  = '%s'`, gov.TagActionDeposit))
	pubsub.Subscribe(queryDeposit, gov.ReadyForVoteHandler)

	pubsub.Subscribe(tmtypes.EventQueryNewBlock, NewBlockEventHandler)
}
