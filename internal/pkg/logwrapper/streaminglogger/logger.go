package streaminglogger

import (
	"github.com/gstdl/crypto-arbitrage/internal/config"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/logwrapperbase"
	"github.com/sirupsen/logrus"
)

type StreamingLogger struct {
	*logwrapperbase.LoggerBase
}

func New(packagePath, exchange string) (logger *StreamingLogger) {

	env := config.GetEnvironmentValue()

	logger = &StreamingLogger{logwrapperbase.NewLogger(packagePath, exchange, "streaming")}

	if env == "development" {
		logger.LogToFileWithLevel("logs/streaming_logs.log", []logrus.Level{logrus.InfoLevel})
	} else {
		logger.LogToFileWithLevel("logs/streaming_logs.log", []logrus.Level{logrus.WarnLevel})
	}

	return
}
