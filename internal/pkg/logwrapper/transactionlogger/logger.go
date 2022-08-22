package transactionlogger

import (
	"github.com/gstdl/crypto-arbitrage/internal/config"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/logwrapperbase"
)

type TransactionLogger struct {
	*logwrapperbase.LoggerBase
}

func New(packagePath, exchange string) (logger *TransactionLogger) {

	logger = &TransactionLogger{logwrapperbase.NewLogger(packagePath, exchange, "transactions")}

	if config.GetEnvironmentValue() == "production" {
		logger.LogToFile("logs/transaction_logs.log")
	} else {
		logger.LogToFile("logs/transaction_forward_testing_logs.log")
	}

	return
}
