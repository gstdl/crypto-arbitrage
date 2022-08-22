package luno

import (
	"context"
	"time"

	"github.com/gstdl/crypto-arbitrage/internal/config"
	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

// func (p *PreOrder) SetClient(client LunoClient) {
// 	p.Client = client.client
// }

func (p PreOrder) GetPairName() (pairName string) {
	pairName = p.pairName
	return
}

func (p PreOrder) GetOrderType() (orderType string) {
	orderType = string(p.orderType)
	return
}

func (p PreOrder) OrderBase() (base string) {
	base = p.base
	return
}

func (p PreOrder) OrderResult() (result string) {
	result = p.result
	return
}

func (p PreOrder) GetThreshold() (threshold decimal.Decimal) {
	threshold, err := p.Client.assets.GetThreshold(p.base)
	if err != nil {
		panic(err)
	}
	return
}

func (p PreOrder) SimulateOrder(streamedPrice interface{}, initialBalance decimal.Decimal) (newBalance decimal.Decimal, ok bool) {

	streamResult, ok2 := streamedPrice.(StreamMap)
	if !ok2 {
		return
	}

	adjustedBalance := round(p.base, initialBalance)

	if streamResult[p.pairName].IsActive() {
		switch p.orderType {
		case luno.OrderTypeBuy:
			newBalance = adjustedBalance.Div(streamResult[p.pairName].newAskPrice, 64)
			ok = newBalance.Sub(p.minVolume).String()[0] != '-'
		case luno.OrderTypeSell:
			ok = adjustedBalance.Sub(p.minVolume).String()[0] != '-'
			if ok {
				newBalance = adjustedBalance.Mul(streamResult[p.pairName].newBidPrice)
			}
		}
	}
	// fmt.Printf("--%v: %v\n%v: %v -> %v ->> %v \n\n", p.pairName, p.orderType, p.base, initialBalance.String(), adjustedBalance.String(), newBalance.String())

	return
}

func (p PreOrder) makeMarketOrderRequest(tradeAmount decimal.Decimal) (req luno.PostMarketOrderRequest) {
	assetAddreses := p.Client.GetSpotWalletAddress()

	if p.orderType == luno.OrderTypeBuy {
		baseId, counterId := assetAddreses[p.result], assetAddreses[p.base]
		tradeAmount = round(p.base, tradeAmount)
		req = luno.PostMarketOrderRequest{
			Pair:             p.pairName,
			Type:             luno.OrderTypeBuy,
			BaseAccountId:    baseId,
			CounterAccountId: counterId,
			CounterVolume:    tradeAmount,
		}
	} else {
		counterId, baseId := assetAddreses[p.result], assetAddreses[p.base]
		tradeAmount = round(p.base, tradeAmount)
		req = luno.PostMarketOrderRequest{
			Pair:             p.pairName,
			Type:             luno.OrderTypeSell,
			BaseAccountId:    baseId,
			BaseVolume:       tradeAmount,
			CounterAccountId: counterId,
		}
	}
	return
}
func (p PreOrder) ExecuteMarketOrderRequest(streamedPrice interface{}, tradeAmount decimal.Decimal, recordId time.Time, index int) (newBalance decimal.Decimal, err error) {

	req := p.makeMarketOrderRequest(tradeAmount)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	var response *luno.PostMarketOrderResponse
	if config.GetEnvironmentValue() == "production" {
		response, err = p.Client.client.PostMarketOrder(ctx, &req)
		if err != nil {
			p.Client.StandardLogger.LogRequestErrors(req, err)
			return
		}
	}

	// go func(tradeAmount decimal.Decimal, index int, streamedPrice, req, response interface{}) {
	newBalance, _ = p.SimulateOrder(streamedPrice, tradeAmount)
	p.Client.TransactionLogger.LogSingleTransaction(p.pairName, p.orderType, tradeAmount.String(), newBalance.String(), "ACTIVE", req, response, recordId, index)
	// }(tradeAmount, index, streamedPrice, req, response)

	return
}
