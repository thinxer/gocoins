package gocoins

import (
	"log"
	"time"
)

// This function is used for tailing trading data of a Client
// by repeatly calling History function.
// Trades are returned to the chan Trade t.
// Some APIs (like BTC-E) use timestamp for since parameter.
// You should set useTimestamp to true if that is the case.
func Tail(c Client, pair Pair, interval time.Duration, useTimestamp bool, t chan Trade) {
	var tid, timestamp int64
	tid = -1
	timestamp = -1
	for {
		start := time.Now()
		since := tid
		if useTimestamp {
			since = timestamp
		}
		trades, err := c.History(pair, since)
		dur := time.Now().Sub(start)
		if err == nil {
			for _, tx := range trades {
				if tx.Id > tid {
					t <- tx
				}
			}
			if len(trades) > 0 {
				tid = trades[len(trades)-1].Id
				timestamp = trades[len(trades)-1].Timestamp
			}
		} else {
			log.Print(err.Error())
			close(t)
			break
		}
		time.Sleep(interval - dur)
	}
}