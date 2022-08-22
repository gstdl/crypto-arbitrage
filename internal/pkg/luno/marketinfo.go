package luno

import (
	"context"
	"strings"
	"time"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func (lc *LunoClient) GetTickers() (tickers []luno.Ticker) {
	req := luno.GetTickersRequest{}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	res, err := lc.client.GetTickers(ctx, &req)
	if err != nil {
		lc.StandardLogger.ErrorFunc("*LunoClient.GetTickers", err)
	}

	tickers = res.Tickers
	return
}

func (lc *LunoClient) GetTicker(pair string) (res *luno.GetTickerResponse) {
	req := luno.GetTickerRequest{Pair: pair}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	res, err := lc.client.GetTicker(ctx, &req)
	if err != nil {
		lc.StandardLogger.ErrorFunc("*LunoClient.GetTicker", err)
	}
	return
}

func (lc *LunoClient) GetFeeInfo(pair string) (res *luno.GetFeeInfoResponse) {
	req := luno.GetFeeInfoRequest{Pair: pair}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	res, err := lc.client.GetFeeInfo(ctx, &req)
	if err != nil {
		lc.StandardLogger.ErrorFunc("*LunoClient.GetFeeInfo", err)
	}
	return
}

func (lc *LunoClient) GetMarketInfo() (info []luno.MarketInfo) {
	req := luno.MarketsRequest{}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	res, err := lc.client.Markets(ctx, &req)
	if err != nil {
		lc.StandardLogger.ErrorFunc("*LunoClient.GetMarketInfo", err)
	}
	info = res.Markets
	return
}

func (lc *LunoClient) GetPairInfo() (setup map[string][]PreOrder, pairs map[string]bool) {
	setup = make(map[string][]PreOrder)
	pairs = make(map[string]bool)

	marketInfo := lc.GetMarketInfo()
	assets := *lc.assets

	for _, info := range marketInfo {
		_, okBase := assets[strings.ToLower(info.BaseCurrency)]
		_, okCounter := assets[strings.ToLower(info.CounterCurrency)]
		if okBase && okCounter {
			pairs[info.MarketId] = true

			fee := lc.GetFeeInfo(info.MarketId)

			if takerFee, err := decimal.NewFromString(fee.TakerFee); err == nil {

				if path, ok := setup[info.CounterCurrency]; ok {
					setup[info.CounterCurrency] = append([]PreOrder{}, append(path, PreOrder{
						pairName:  info.MarketId,
						base:      info.CounterCurrency,
						result:    info.BaseCurrency,
						orderType: luno.OrderTypeBuy,
						takerFee:  takerFee,
						minVolume: info.MinVolume,
						Client:    lc,
					})...)
				} else {
					setup[info.CounterCurrency] = []PreOrder{{
						pairName:  info.MarketId,
						base:      info.CounterCurrency,
						result:    info.BaseCurrency,
						orderType: luno.OrderTypeBuy,
						takerFee:  takerFee,
						minVolume: info.MinVolume,
						Client:    lc,
					}}
				}

				if path, ok := setup[info.BaseCurrency]; ok {
					setup[info.BaseCurrency] = append([]PreOrder{}, append(path, PreOrder{
						pairName:  info.MarketId,
						base:      info.BaseCurrency,
						result:    info.CounterCurrency,
						orderType: luno.OrderTypeSell,
						takerFee:  takerFee,
						minVolume: info.MinVolume,
						Client:    lc,
					})...)
				} else {
					setup[info.BaseCurrency] = []PreOrder{{
						pairName:  info.MarketId,
						base:      info.BaseCurrency,
						result:    info.CounterCurrency,
						orderType: luno.OrderTypeSell,
						takerFee:  takerFee,
						minVolume: info.MinVolume,
						Client:    lc,
					}}
				}
			}
		}
	}
	return
}
