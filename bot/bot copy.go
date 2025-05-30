package bot

// processCandle runs the strategy logic on each closed candle
// func processCandleOri(c Candle, symbol, token, interval string) {
// 	s := strings.ToUpper(symbol[:len(symbol)-4])

// 	closes = append(closes, c.Close)
// 	volumes = append(volumes, c.Volume)

// 	if len(closes) > 500 {
// 		closes = closes[1:]
// 		volumes = volumes[1:]
// 	}

// 	if len(closes) < rsiLength || len(volumes) < volumeLookback || len(closes) < bbLength {
// 		return
// 	}

// 	step := constant.StepMap[symbol]

// 	pricePrecision := constant.SymbolPrecisionMap[symbol][0]
// 	amountPrecision := constant.SymbolPrecisionMap[symbol][1]

// 	price := fmt.Sprintf("%%.%df", pricePrecision)
// 	amount := fmt.Sprintf("%%.%df", amountPrecision)

// 	rsiVal := calcRSI(closes, rsiLength)
// 	avgVolume := sma(volumes, volumeLookback)
// 	highVolume := c.Volume > avgVolume*1.5
// 	extremeHighVolume := c.Volume > avgVolume*2

// 	greenCandle := c.Close > c.Open
// 	redCandle := c.Close < c.Open

// 	highLowDiff := c.High - c.Low
// 	if highLowDiff == 0 {
// 		return
// 	}
// 	topWick := c.High - math.Max(c.Open, c.Close)
// 	bottomWick := math.Min(c.Open, c.Close) - c.Low
// 	topWickPerc := (topWick / highLowDiff) * 100
// 	bottomWickPerc := (bottomWick / highLowDiff) * 100

// 	basis := sma(closes[len(closes)-bbLength:], bbLength)
// 	stdDev := stddev(closes[len(closes)-bbLength:], basis)
// 	upper := basis + bbMult*stdDev
// 	lower := basis - bbMult*stdDev

// 	// sudden spike up/down
// 	spikeUpPerc := ((c.High - c.Open) / c.Open * 100)
// 	spikeDownPerc := ((c.Low - c.Open) / c.Open * 100)
// 	percent := constant.PercentageMap[interval]

// 	// long needle and close as red candle with volumex2
// 	wickSpikeUp := spikeUpPerc >= percent && redCandle && extremeHighVolume
// 	// long needle and close as green candle with volumex2
// 	wickSpikeDown := spikeDownPerc <= percent && greenCandle && extremeHighVolume

// 	rawBuy := (rsiVal < 35 && highVolume && (greenCandle || (redCandle && bottomWickPerc > 60))) || wickSpikeDown
// 	rawSell := (rsiVal > 65 && highVolume && (redCandle || (greenCandle && topWickPerc > 60))) || wickSpikeUp

// 	combinedBuy := rawBuy && c.Close <= lower
// 	combinedSell := rawSell && c.Close >= upper

// 	buySignal := combinedBuy
// 	sellSignal := combinedSell

// 	// === STOP LOSS CHECK ===
// 	if state == 1 && c.Close <= entryPrice*(1-stopLossPercent/100) {
// 		profit := (c.Close - entryPrice) * positionSize
// 		percentChange := ((c.Close - entryPrice) / entryPrice) * 100
// 		balance += tradeUSDT + profit
// 		// placeOrder(symbol, "SELL", positionSize)
// 		a := fmt.Sprintf("STOP LOSS [LONG]\nAmount: "+amount+"%s"+"\nPrice: "+price+"\nPercent changed: %.2f\nLoss: %.2f USDT\nBalance: %.2f USDT\n", positionSize, s, c.Close, percentChange, profit, balance)
// 		log.Println(a)
// 		sendTelegramMessage(token, a)
// 		state, positionSize, entryPrice = 0, 0, 0
// 		totalProfitLoss += profit
// 		b := fmt.Sprintf("Total profit/loss : %.2f", totalProfitLoss)
// 		log.Println(b)
// 		sendTelegramMessage(token, b)
// 		return
// 	}
// 	if state == -1 && c.Close >= entryPrice*(1+stopLossPercent/100) {
// 		closeAmount := math.Abs(positionSize)
// 		profit := (entryPrice - c.Close) * closeAmount
// 		percentChange := ((c.Close - entryPrice) / entryPrice) * 100
// 		balance += tradeUSDT + profit
// 		// placeOrder(symbol, "BUY", positionSize)
// 		a := fmt.Sprintf("STOP LOSS [SHORT]\nAmount: "+amount+"%s"+"\nPrice: "+price+"\nPercent changed: %.2f\nLoss: %.2f USDT\nBalance: %.2f USDT\n", closeAmount, s, c.Close, percentChange, profit, balance)
// 		log.Println(a)
// 		sendTelegramMessage(token, a)
// 		state, positionSize, entryPrice = 0, 0, 0
// 		totalProfitLoss += profit
// 		b := fmt.Sprintf("Total profit/loss : %.2f", totalProfitLoss)
// 		log.Println(b)
// 		sendTelegramMessage(token, b)
// 		return
// 	}

// 	// === TRADING LOGIC ===
// 	if state == 0 {
// 		// Neutral: open position on any signal
// 		if buySignal {
// 			if balance >= tradeUSDT {
// 				size := roundUpToStep(tradeUSDT/c.Close, step)
// 				positionSize = size
// 				entryPrice = c.Close
// 				balance -= tradeUSDT
// 				state = 1
// 				// placeOrder(symbol, "BUY", positionSize)
// 				a := fmt.Sprintf("[LONG]\nAmount: "+amount+"%s"+"\nPrice: "+price+"\nStop loss: "+price+"\nBalance: %.2f", size, s, c.Close, c.Close*(1-stopLossPercent/100), balance)
// 				log.Println(a)
// 				sendTelegramMessage(token, a)
// 			} else {
// 				a := "Insufficient balance to open LONG position"
// 				log.Println(a)
// 				sendTelegramMessage(token, a)
// 			}
// 			return
// 		}
// 		if sellSignal {
// 			if balance >= tradeUSDT {
// 				size := roundUpToStep(tradeUSDT/c.Close, step)
// 				positionSize = -size
// 				entryPrice = c.Close
// 				balance -= tradeUSDT
// 				state = -1
// 				// placeOrder(symbol, "SELL", positionSize)
// 				a := fmt.Sprintf("[SHORT]\nAmount: "+amount+"%s"+"\nPrice: "+price+"\nStop loss: "+price+"\nBalance: %.2f", size, s, c.Close, c.Close*(1+stopLossPercent/100), balance)
// 				log.Println(a)
// 				sendTelegramMessage(token, a)
// 			} else {
// 				a := "Insufficient balance to open SHORT position"
// 				log.Println(a)
// 				sendTelegramMessage(token, a)
// 			}
// 			return
// 		}
// 	} else if state == 1 {
// 		// Long position: close only on sell signal
// 		if sellSignal {
// 			profit := (c.Close - entryPrice) * positionSize
// 			percentChange := ((c.Close - entryPrice) / entryPrice) * 100
// 			balance += tradeUSDT + profit
// 			// placeOrder(symbol, "SELL", positionSize)
// 			a := fmt.Sprintf("Closed [LONG]\nAmount: "+amount+"%s"+"\nPrice: "+price+"\nPercent changed: %.2f\nProfit: %.2f USDT\nBalance: %.2f USDT\n", positionSize, s, c.Close, percentChange, profit, balance)
// 			log.Println(a)
// 			sendTelegramMessage(token, a)
// 			state, positionSize, entryPrice = 0, 0, 0
// 			totalProfitLoss += profit
// 			b := fmt.Sprintf("Total profit/loss : %.2f", totalProfitLoss)
// 			log.Println(b)
// 			sendTelegramMessage(token, b)
// 			return
// 		}
// 	} else if state == -1 {
// 		// Short position: close only on buy signal
// 		if buySignal {
// 			closeAmount := math.Abs(positionSize)
// 			profit := (entryPrice - c.Close) * closeAmount
// 			percentChange := ((c.Close - entryPrice) / entryPrice) * 100
// 			balance += tradeUSDT + profit
// 			// placeOrder(symbol, "BUY", positionSize)
// 			a := fmt.Sprintf("Closed [SHORT]\nAmount: "+amount+"%s"+"\nPrice: "+price+"\nPercent changed: %.2f\nProfit: %.2f USDT\nBalance: %.2f USDT\n", closeAmount, s, c.Close, percentChange, profit, balance)
// 			log.Println(a)
// 			sendTelegramMessage(token, a)
// 			state, positionSize, entryPrice = 0, 0, 0
// 			totalProfitLoss += profit
// 			b := fmt.Sprintf("Total profit/loss : %.2f", totalProfitLoss)
// 			log.Println(b)
// 			sendTelegramMessage(token, b)
// 			return
// 		}
// 	}
// }
