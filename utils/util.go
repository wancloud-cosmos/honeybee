package utils

import (
	amino "github.com/tendermint/go-amino"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func init() {
	ctypes.RegisterAmino(CDC)
}

var CDC = amino.NewCodec()
