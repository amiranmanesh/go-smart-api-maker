package main

import (
	"github.com/go-kit/kit/log"
	"os"
)

func main() {
	var accountLogger log.Logger
	{
		accountLogger = log.NewLogfmtLogger(os.Stderr)
		accountLogger = log.NewSyncLogger(accountLogger)
		accountLogger = log.With(accountLogger,
			"service", "account",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
}
