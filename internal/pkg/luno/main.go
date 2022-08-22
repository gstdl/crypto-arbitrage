package luno

import (
	"strings"

	"github.com/gstdl/crypto-arbitrage/internal/config"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/standardlogger"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/transactionlogger"
	"github.com/luno/luno-go"
)

func NewClient() (lc *LunoClient) {
	// Read config
	config := config.GetLunoConfig()

	// Setup luno client
	client := luno.NewClient()
	client.SetAuth(config.Id, config.Secret)

	// Setup standardlogger
	StandardLogger := standardlogger.New("internal/pkg/luno", "LUNO")

	// Setup transactionlogger
	transactionLogger := transactionlogger.New("internal/pkg/luno", "LUNO")

	lc = &LunoClient{
		client: client,
		StandardLogger: StandardLogger,
		TransactionLogger: transactionLogger,
		assets: &config.Assets,
	}

	StandardLogger.LogServiceStart()
	
	return
}

func (lc *LunoClient) GetAssetNames() (assets []string) {
	for asset := range *lc.assets {
		assets = append(assets, strings.ToUpper(asset))
	}
	return
}
