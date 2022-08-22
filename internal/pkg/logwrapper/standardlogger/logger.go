package standardlogger

import (
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/logwrapperbase"
)

type StandardLogger struct {
	*logwrapperbase.LoggerBase
}

func New(packagePath, exchange string) (logger *StandardLogger) {

	// environment := config.GetEnvironmentValue()

	logger = &StandardLogger{logwrapperbase.NewLogger(packagePath, exchange, "standard")}

	logger.LogToFile("logs/standard_logs.log")
	return
}
