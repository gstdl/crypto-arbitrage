package luno

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gstdl/crypto-arbitrage/internal/config"
	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func (lc *LunoClient) GetBalance() (balances []luno.AccountBalance) {
	req := luno.GetBalancesRequest{Assets: lc.GetAssetNames()}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	res, err := lc.client.GetBalances(ctx, &req)
	if err != nil {
		lc.StandardLogger.ErrorFunc("*LunoClient.GetBalance", err)
	}

	balances = res.Balance
	return
}

func (lc *LunoClient) GetWalletBalance() (balance map[string]map[string]decimal.Decimal) {
	balance = make(map[string]map[string]decimal.Decimal)
	for _, key := range [2]string{"SPOT", "SAVINGS"} {
		balance[key] = make(map[string]decimal.Decimal)
	}

	lunoBalance := lc.GetBalance()
	assets := *lc.assets

	for _, assetBalance := range lunoBalance {
		if assetInfo, ok := assets[strings.ToLower(assetBalance.Asset)]; ok {
			if config.GetEnvironmentValue() == "production" {
				if assetInfo.AddressID == assetBalance.AccountId {
					balance["SPOT"][assetBalance.Asset] = assetBalance.Balance
				} else {
					balance["SAVINGS"][assetBalance.Asset] = assetBalance.Balance
				}
			} else {
				var devBalance float64
				switch assetBalance.Asset {
				case "XBT":
					devBalance = 0.15
				case "ETH":
					devBalance = 1.0
				case "USDC":
					devBalance = 5000.0
				case "IDR":
					devBalance = 50000000.0
				}
				balance["SPOT"][assetBalance.Asset] = decimal.NewFromFloat64(devBalance, 6)
				balance["SAVINGS"][assetBalance.Asset] = decimal.NewFromFloat64(devBalance, 6)
			}
		}
	}

	// lc.StandardLogger.Logger.Info("Successfully retrieved balance")

	return
}

func (lc *LunoClient) GetSpotWalletBalance() (balance map[string]decimal.Decimal) {
	balance = lc.GetWalletBalance()["SPOT"]
	return
}

func (lc *LunoClient) GetSavingsWalletBalance() (balance map[string]decimal.Decimal) {
	balance = lc.GetWalletBalance()["SAVINGS"]
	return
}

func (lc *LunoClient) GetSpotWalletAddress() (addresses map[string]int64) {
	addresses = make(map[string]int64)

	assets := *lc.assets

	for currency, values := range assets {
		addressStr := values.AddressID
		address, _ := strconv.ParseInt(addressStr, 8, 64)
		addresses[currency] = address
	}

	return

}
