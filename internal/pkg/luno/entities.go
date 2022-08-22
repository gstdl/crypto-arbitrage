package luno

import (
	"github.com/gstdl/crypto-arbitrage/internal/config"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/standardlogger"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/streaminglogger"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/transactionlogger"
	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"github.com/luno/luno-go/streaming"
)

type LunoClient struct {
	client            *luno.Client
	StandardLogger    *standardlogger.StandardLogger
	TransactionLogger *transactionlogger.TransactionLogger
	assets            *config.LunoAssetMap
}

type PreOrder struct {
	pairName  string
	base      string
	result    string
	orderType luno.OrderType
	takerFee  decimal.Decimal
	minVolume decimal.Decimal
	Client    *LunoClient
}

type lunoStreamer struct {
	conn     *streaming.Conn
	logger   *streaminglogger.StreamingLogger
	pairName string
}

type streamResult struct {
	newBidPrice  decimal.Decimal
	newAskPrice  decimal.Decimal
	newStatus luno.Status
	PairName  string
}

type StreamMap map[string]streamResult
