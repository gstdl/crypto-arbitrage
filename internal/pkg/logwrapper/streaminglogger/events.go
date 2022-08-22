package streaminglogger

import (
	"strconv"

	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/logwrapperbase"
	"github.com/sirupsen/logrus"
)

func newEvent(logCategory int, logSubCategory int, eventLevel logrus.Level, message string) (e logwrapperbase.Event) {
	e = logwrapperbase.NewEvent(logCategory, logSubCategory, eventLevel, message, "streaminglogger")
	return
}

// Declare variables to store log messages as new Events
var (
	streamingSuccessMessage = newEvent(200, 0, logrus.InfoLevel, "Streaming %s Success!")
)

func (l *StreamingLogger) LogSuccessStream(pairName, askingPriceStr, bidPriceStr, status interface{}) {
	go func() {
		askingPrice, _ := strconv.ParseFloat(askingPriceStr.(string), 64)
		bidPrice, _ := strconv.ParseFloat(bidPriceStr.(string), 64)
		fields := map[string]interface{}{
			"pair":      pairName,
			"ask_price": askingPrice,
			"bid_price": bidPrice,
			"status":    status,
		}
		l.LogWithFields(streamingSuccessMessage, fields, pairName)
	}()
}
