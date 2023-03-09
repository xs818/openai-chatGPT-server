package util

import (
	logging "github.com/ipfs/go-log/v2"
	"os"
)

var Logger = logging.Logger("Server log")

func init() {
	lvl, set := os.LookupEnv("GOLOG_LOG_LEVEL")
	if !set {
		lvl, err := logging.LevelFromString("debug")
		if err != nil {
			panic(err)
		}
		logging.SetAllLoggers(lvl)
		return
	}

	_ = logging.SetLogLevel("*", lvl)

}
