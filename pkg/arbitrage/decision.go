package arbitrage

import (
	"sync"
	"time"

	"github.com/luno/luno-go/decimal"
)

func (as ArbitrageService) MakeDecision(balance map[string]decimal.Decimal, streamedPrice interface{}) {
	exectuableTrades := make(chan executable)
	wg := sync.WaitGroup{}

	pm := as.paths

	for starterCurrency, paths := range pm {
		initialBalance := balance[starterCurrency]
		for _, path := range paths {
			wg.Add(1)
			go func(p Path) {
				defer wg.Done()
				newBalance, pairNames, isProfitable, ok := p.SimulateOrders(initialBalance, streamedPrice)
				if isProfitable && ok {
					recordedTime := time.Now()
					exectuableTrades <- executable{
						currency:       p[0].OrderBase(),
						initialBalance: initialBalance,
						newBalance:     newBalance,
						prices:         streamedPrice,
						pairHistory:    pairNames,
						path:           p,
						recordedTime:   recordedTime,
					}
					go func() {
						diff := newBalance.Sub(initialBalance).ToScale(4).String()
						as.logger.StartingArbitrage(diff, pairNames, recordedTime)
					}()
				}
			}(path)
		}
	}

	go func() {
		wg.Wait()
		close(exectuableTrades)
		// fmt.Println("channel closed!")
	}()

	for ex := range exectuableTrades {
		as.ExecuteOrders(ex)
	}
	// fmt.Println("reopening channel!")

}

func (p Path) SimulateOrders(initialBalance decimal.Decimal, streamedPrice interface{}) (newBalance decimal.Decimal, pairNames [][2]string, isProfitable bool, ok bool) {
	var threshold decimal.Decimal

	ok = true
	for ix, preOrder := range p {
		pairNames = append(pairNames, [2]string{preOrder.GetPairName(), preOrder.GetOrderType()})
		if !ok {
			// fmt.Printf("NOT OK! pairNames: %v\n", pairNames)
			return
		}
		if ix == 0 {
			newBalance, ok = preOrder.SimulateOrder(streamedPrice, initialBalance)
			threshold = preOrder.GetThreshold()
		} else {
			// fmt.Printf("newBalance: %v - %v -- %v\n", pairNames, preOrder.GetPairName(), newBalance)
			newBalance, ok = preOrder.SimulateOrder(streamedPrice, newBalance)
		}
	}
	diffStr := newBalance.Sub(initialBalance).Sub(threshold).String()
	// fmt.Printf("diffStr: %v\npairNames: %v\nthresh: %v\n\n", diffStr, pairNames, threshold)
	isProfitable = diffStr[0] != '-'
	return
}

func (as ArbitrageService) ExecuteOrders(ex executable) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ex.execute()
	}()
	// go func() {
	wg.Wait()
	as.logger.LogArbitrageTransaction(ex.currency, ex.pairHistory, ex.initialBalance.String(), ex.newBalance.String(), ex.recordedTime)
	// }()
}

func (ex executable) execute() {
	p := ex.path
	var newBalance decimal.Decimal
	for ix, preOrder := range p {
		if ix == 0 {
			newBalance, _ = preOrder.ExecuteMarketOrderRequest(ex.prices, ex.initialBalance, ex.recordedTime, ix)
		} else {
			newBalance, _ = preOrder.ExecuteMarketOrderRequest(ex.prices, newBalance, ex.recordedTime, ix)
		}
	}
}
