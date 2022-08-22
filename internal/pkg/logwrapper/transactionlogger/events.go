package transactionlogger

import (
	"strconv"

	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/logwrapperbase"
	"github.com/sirupsen/logrus"
)

func newEvent(logCategory int, logSubCategory int, eventLevel logrus.Level, message string) (e logwrapperbase.Event) {
	e = logwrapperbase.NewEvent(logCategory, logSubCategory, eventLevel, message, "transactionlogger")
	return
}

// Declare variables to store log messages as new Events
var (
	arbitrageStartMessage   = newEvent(100, 0, logrus.InfoLevel, "Found arbit opportunity! Expected profit: %v")
	tradeSuccessMessage     = newEvent(200, 0, logrus.InfoLevel, "Transaction Success!")
	arbitrageSuccessMessage = newEvent(200, 0, logrus.InfoLevel, "Arbit %v Success!")
)

func (l *TransactionLogger) StartingArbitrage(expectedProfit, pairs, recordedId interface{}) {
	fields := map[string]interface{}{
		"recordId": recordedId,
		"pairs":    pairs,
	}
	l.LogWithFields(arbitrageStartMessage, fields, expectedProfit)
}

func (l *TransactionLogger) LogSingleTransaction(pairName, orderType, baseBalanceStr, resultBalanceStr, status, request, response, recordedId, recordIndex interface{}) {
	baseBalance, _ := strconv.ParseFloat(baseBalanceStr.(string), 64)
	resultBalance, _ := strconv.ParseFloat(resultBalanceStr.(string), 64)
	fields := map[string]interface{}{
		"type":              "single transation",
		"pair":              pairName,
		"order_type":        orderType,
		"base_balance":      baseBalance,
		"resultant_balance": resultBalance,
		"status":            status,
		"response":          response,
		"recordId":          recordedId,
		"recordIndex":       recordIndex,
	}
	l.LogWithFields(tradeSuccessMessage, fields)
}

func (l *TransactionLogger) LogArbitrageTransaction(base, path, initialBalanceStr, newBalanceStr, recordedId interface{}) {
	initialBalance, _ := strconv.ParseFloat(initialBalanceStr.(string), 64)
	newBalance, _ := strconv.ParseFloat(newBalanceStr.(string), 64)
	fields := map[string]interface{}{
		"type":            "arbit transation",
		"base":            base,
		"path":            path,
		"initial_balance": initialBalance,
		"new_balance":     newBalance,
		"recordId":        recordedId,
	}
	l.LogWithFields(arbitrageSuccessMessage, fields, base)
}
