package main

import (
	"fmt"
	"time"

	"github.com/gstdl/crypto-arbitrage/internal/pkg/luno"
	"github.com/gstdl/crypto-arbitrage/pkg/arbitrage"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			main()
		}
	}()

	printServiceLogo()

	var client = luno.NewClient()

	pairSetup, pairs := client.GetPairInfo()

	pathMap := arbitrage.FindPathsFromLunoPairSetup(pairSetup)

	priceChannel := make(chan luno.StreamMap)
	errorChannel := make(chan interface{}, len(pairSetup))

	go func() {
		luno.StartStreaming(client, pairs, time.Second*30, priceChannel, errorChannel)
	}()

	for {
		select {
		case prices := <-priceChannel:
			// get balances
			balance := client.GetSpotWalletBalance()
			// make decision
			pathMap.MakeDecision(balance, prices)

		case err := <-errorChannel:
			// restart app when it fails
			client.StandardLogger.LogServiceStopped(err)
			panic(err)
		}

	}
}

func printServiceLogo() {
	fmt.Println(`                             __                      _____       ___.   .__  __                                
  ___________ ___.__._______/  |_  ____             /  _  \______\_ |__ |__|/  |_____________     ____   ____  
_/ ___\_  __ <   |  |\____ \   __\/  _ \   ______  /  /_\  \_  __ \ __ \|  \   __\_  __ \__  \   / ___\_/ __ \ 
\  \___|  | \/\___  ||  |_> >  | (  <_> ) /_____/ /    |    \  | \/ \_\ \  ||  |  |  | \// __ \_/ /_/  >  ___/ 
 \___  >__|   / ____||   __/|__|  \____/          \____|__  /__|  |___  /__||__|  |__|  (____  /\___  / \___  >
     \/       \/     |__|                                 \/          \/                     \//_____/      \/ 
                          .__                                                                                  
  ______ ______________  _|__| ____  ____                                                                      
 /  ___// __ \_  __ \  \/ /  |/ ___\/ __ \                                                                     
 \___ \\  ___/|  | \/\   /|  \  \__\  ___/                                                                     
/____  >\___  >__|    \_/ |__|\___  >___  >                                                                    
     \/     \/                    \/    \/  `)
}
