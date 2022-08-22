package luno

import (
	"context"
	"strconv"
	"time"

	"github.com/luno/luno-go"
)

func (lc *LunoClient) GetUserTradeHistory(pair string) (trades []luno.Trade) {
	req := luno.ListUserTradesRequest{Pair: pair}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	res, err := lc.client.ListUserTrades(ctx, &req)
	if err != nil {
		lc.StandardLogger.ErrorFunc("*LunoClient.GetUserTradeHistory", err)
	}
	trades = res.Trades
	return
}

func (lc *LunoClient) GetTransactions() (trades []luno.Transaction) {
	addressID, _ := strconv.Atoi((*lc.assets)["eth"].AddressID)
	req := luno.ListTransactionsRequest{Id: int64(addressID)}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	res, err := lc.client.ListTransactions(ctx, &req)
	if err != nil {
		lc.StandardLogger.ErrorFunc("*LunoClient.GetTransactions", err)
	}
	lc.StandardLogger.Printf("%+v", res)
	trades = res.Transactions
	return
}
