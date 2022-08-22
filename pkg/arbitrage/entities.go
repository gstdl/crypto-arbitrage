package arbitrage

import (
	"time"

	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/transactionlogger"
	"github.com/luno/luno-go/decimal"
)

type PreOrder interface {
	GetPairName() string
	GetOrderType() string
	OrderBase() string
	OrderResult() string
	GetThreshold() decimal.Decimal
	SimulateOrder(interface{}, decimal.Decimal) (decimal.Decimal, bool)
	ExecuteMarketOrderRequest(interface{}, decimal.Decimal, time.Time, int) (decimal.Decimal, error)
}

type ArbitrageService struct {
	paths  PathMap
	logger *transactionlogger.TransactionLogger
}

type executable struct {
	currency string
	initialBalance  decimal.Decimal
	newBalance      decimal.Decimal
	pairHistory     [][2]string
	path            Path
	prices          interface{}
	recordedTime    time.Time
}

type PathMap map[string][]Path

type Path []PreOrder

type queue struct {
	elements []interface{}
}
