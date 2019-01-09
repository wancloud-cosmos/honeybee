package utils

import (
	"github.com/astaxie/beego"

	amino "github.com/tendermint/go-amino"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func init() {
	ctypes.RegisterAmino(CDC)
}

var CDC = amino.NewCodec()

func LogLevel() int {
	levelStr := beego.AppConfig.String("log::level")
	if "" == levelStr {
		return beego.LevelDebug
	}

	switch levelStr {
	case "debug":
		return beego.LevelDebug
	case "info":
		return beego.LevelInformational
	case "warning":
		return beego.LevelWarning
	case "error":
		return beego.LevelError
	}

	return beego.LevelDebug
}
