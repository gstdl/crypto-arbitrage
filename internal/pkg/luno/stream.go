package luno

import (
	"time"

	"github.com/gstdl/crypto-arbitrage/internal/config"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/streaminglogger"
	lunogo "github.com/luno/luno-go"
	"github.com/luno/luno-go/streaming"
)

// Setup logwrapper
var logger = streaminglogger.New("internal/pkg/luno-stream", "LUNO")

func StartStreaming(lc *LunoClient, pairs map[string]bool, streamDelay time.Duration, listenerChannel chan StreamMap, errorChannel chan interface{}) {

	// Read config
	id, secret := config.GetLunoCredentials()

	streamChan := make(chan streamResult)

	go func() {
		for pair := range pairs {
			go func(pair string) {
				defer func() {
					if r := recover(); r != nil {
						errorChannel <- r
					}
				}()
				// Setup luno streaming
				c, err := streaming.Dial(id, secret, pair)
				if err != nil {
					logger.Fatal(err)
				}
				defer c.Close()

				ls := lunoStreamer{
					conn:     c,
					logger:   logger,
					pairName: pair,
				}

				for {
					time.Sleep(streamDelay)
					ls.logSnapshot(streamChan)
				}
			}(pair)
		}
	}()

	event := make(StreamMap)
	for {
		e := <-streamChan
		event[e.PairName] = e
		if len(event) == 5 {

			listenerChannel <- event

			event = make(StreamMap)

		}

	}
}

func (ls *lunoStreamer) logSnapshot(listenerChannel chan streamResult) {
	ss := ls.conn.Snapshot()
	askingPrice := ss.Asks[0].Price
	bidPrice := ss.Bids[0].Price
	listenerChannel <- streamResult{
		newAskPrice: askingPrice,
		newBidPrice: bidPrice,
		newStatus:   ss.Status,
		PairName:    ls.pairName,
	}
	if config.GetEnvironmentValue() == "development" {
		ls.logger.LogSuccessStream(ls.pairName, askingPrice.String(), bidPrice.String(), ss.Status)
	}
}

func (sr streamResult) IsActive() (ok bool) {
	ok = sr.newStatus == lunogo.StatusActive
	return
}
