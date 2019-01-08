package pubsub

import (
	"encoding/json"

	"context"
	"validator-monitor/app/config"
	"validator-monitor/utils"

	"github.com/astaxie/beego"
	"github.com/tendermint/tendermint/libs/events"
	tmpubsub "github.com/tendermint/tendermint/libs/pubsub"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	types "github.com/tendermint/tendermint/rpc/lib/types"

	// "github.com/tendermint/tendermint/libs/pubsub/query"
	cli "github.com/tendermint/tendermint/rpc/lib/client"
)

type SubscribeCallbackFunc func(query string, data events.EventData) error

func EventHandler(respCh chan types.RPCResponse, handler SubscribeCallbackFunc) {
	for v := range respCh {
		q, data, err := UnmarshalEvent(v.Result)
		if nil != err {
			beego.Error(err)
			return
		}
		beego.Debug("received event:", q, data)

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

func Subscribe(q tmpubsub.Query, cb SubscribeCallbackFunc) (client *cli.WSClient, err error) {
	client = cli.NewWSClient(config.NodeAddr, "/websocket")
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
