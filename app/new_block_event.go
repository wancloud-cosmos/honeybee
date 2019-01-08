package app

import (
	"fmt"
	"validator-monitor/app/config"

	"github.com/astaxie/beego"
	"github.com/tendermint/tendermint/libs/events"
	tmtypes "github.com/tendermint/tendermint/types"
)

func NewBlockEventHandler(query string, data events.EventData) error {
	block, ok := data.(tmtypes.EventDataNewBlock)
	if !ok {
		err := fmt.Errorf("EventData to EventDataNewBlock, conversion failed")
		beego.Error(err.Error())

		//will exit the sub when return error
		return nil
	}

	addr := config.ValidatorAddresses[0]
	commits := block.Block.LastCommit.Precommits
	if !monitorNode.IsInLastCommit(addr, commits) {
		err := fmt.Errorf("addr(%s) miss block(%d)", addr, block.Block.Height-1)
		beego.Error(err)

		defaultMissBlocks.SetMiss(block.Block.Height - 1)

		//will exit the sub when return error
		return nil
	}

	return nil
}
