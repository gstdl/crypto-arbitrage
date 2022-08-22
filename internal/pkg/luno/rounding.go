package luno

import (
	"github.com/luno/luno-go/decimal"
)

func round(currency string, balance decimal.Decimal) (adjustedBalance decimal.Decimal) {
	if currency == "IDR" {
		adjustedBalance = balance.ToScale(-3)
	} else if currency == "XBT" || currency == "BTC" {
		adjustedBalance = balance.ToScale(4)
	} else {
		adjustedBalance = balance.ToScale(2)
	}

	return
}
