package app

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/tendermint/tendermint/libs/events"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	types "github.com/tendermint/tendermint/rpc/lib/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"context"
	"validator-monitor/utils"

	"github.com/tendermint/tendermint/libs/pubsub/query"
	cli "github.com/tendermint/tendermint/rpc/lib/client"
)

func ReadyForVoteHandler(query string, data events.EventData) error {
	tags := data.(tmtypes.EventDataTx).TxResult.Result.Tags
	for _, v := range tags {
		beego.Debug(v.String())
		if "voting-period-start" == string(v.Key) {
			var id int64
			err := utils.CDC.UnmarshalBinaryBare(v.Value, &id)
			if nil != err {
				beego.Error(err)
				return err
			}

			beego.Info("proposal-id:", id, "ready for vote")

			AutoVote(id)

			return nil
		}
	}

	return nil
}

func EventHandler(respCh chan types.RPCResponse, handler SubscribeCallbackFunc) {
	for v := range respCh {
		q, data, err := UnmarshalEvent(v.Result)
		if nil != err {
			beego.Error(err)
			return
		}
		beego.Debug(q, data)

		//TODO should more research to find out why first event is empty
		if "" == q {
			continue
		}

		err = handler(q, data)
		if nil != err {
			beego.Error(err)
			return
		}
	}
}

type SubscribeCallbackFunc func(query string, data events.EventData) error

func Subscribe(q *query.Query, cb SubscribeCallbackFunc) (client *cli.WSClient, err error) {
	client = cli.NewWSClient(monitorNodes[0].Addr, "/websocket")
	err = client.Start()
	if nil != err {
		beego.Error("WSClient start failed,", err)
		return nil, err
	}

	//TODO maybe shoud set timeout
	err = client.Subscribe(context.TODO(), q.String())
	if nil != err {
		beego.Error(err)
		return nil, err
	}

	go EventHandler(client.ResponsesCh, cb)

	return client, nil
}

// UnmarshalEvent unmarshals a json event
func UnmarshalEvent(b json.RawMessage) (string, events.EventData, error) {
	event := new(ctypes.ResultEvent)
	if err := utils.CDC.UnmarshalJSON(b, event); err != nil {
		return "", nil, err
	}
	return event.Query, event.Data, nil
}
