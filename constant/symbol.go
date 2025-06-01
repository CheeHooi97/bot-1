package constant

import "bot-1/config"

func GetTokenMap() map[string]any {
	return map[string]any{
		"btcusdt": map[string]string{
			"1m":  config.BTC1m,
			"5m":  config.BTC5m,
			"15m": config.BTC15m,
			"1h":  config.BTC1h,
			"4h":  config.BTC4h,
			// "1d":  config.BTC1d,
		},
		"ethusdt": map[string]string{
			"1m":  config.ETH1m,
			"5m":  config.ETH5m,
			"15m": config.ETH15m,
			"1h":  config.ETH1h,
			"4h":  config.ETH4h,
			// "1d":  config.ETH1d,
		},
		"solusdt": map[string]string{
			"1m": config.SOL1m,
			// "5m":  config.SOL5m,
			"15m": config.SOL15m,
			// "1h":  config.SOL1h,
			// "4h":  config.SOL4h,
			// "1d":  config.SOL1d,
		},
		"bnbusdt": map[string]string{
			"1m": config.BNB1m,
			// "5m":  config.BNB5m,
			"15m": config.BNB15m,
			// "1h":  config.BNB1h,
			// "4h": config.BNB4h,
			// "1d":  config.BNB1d,
		},
		"adausdt": map[string]string{
			"1m": config.ADA1m,
			// "5m":   config.ADA5m,
			"15m": config.ADA15m,
			// "1h":  config.ADA1h,
			// "4h": config.ADA4h,
			// "1d":  config.ADA1d,
		},
		"xrpusdt": map[string]string{
			"1m": config.XRP1m,
			// "5m":   config.XRP5m,
			"15m": config.XRP15m,
			// "1h":  config.XRP1h,
			// "4h":config.XRP4h,
			// "1d": config.XRP1d,
		},
		"hypeusdt": map[string]string{
			"1m": config.ETH4h,
		},
		"dogeusdt": map[string]string{
			"1m": config.BTC4h,
		},
	}
}
