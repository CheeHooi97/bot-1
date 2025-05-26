package constant

import "bot-1/config"

func GetTokenMap() map[string]any {
	return map[string]any{
		"btcusdt": map[string]string{
			"1m":  config.TelegramTokenBTC1m,
			"5m":  config.TelegramTokenBTC5m,
			"15m": config.TelegramTokenBTC15m,
			"30m": config.TelegramTokenBTC30m,
			"1h":  config.TelegramTokenBTC1h,
			"4h":  config.TelegramTokenBTC4h,
			// "1d":  config.TelegramTokenBTC1d,
		},
		"ethusdt": map[string]string{
			"1m":  config.TelegramTokenETH1m,
			"5m":  config.TelegramTokenETH5m,
			"15m": config.TelegramTokenETH15m,
			"30m": config.TelegramTokenETH30m,
			"1h":  config.TelegramTokenETH1h,
			"4h":  config.TelegramTokenETH4h,
			// "1d":  config.TelegramTokenETH1d,
		},
		"solusdt": map[string]string{
			"1m": config.TelegramTokenSOL1m,
			// "5m":  config.TelegramTokenSOL5m,
			// "15m": config.TelegramTokenSOL15m,
			// "30m": config.TelegramTokenSOL30m,
			// "1h":  config.TelegramTokenSOL1h,
			// "4h":  config.TelegramTokenSOL4h,
			// "1d":  config.TelegramTokenSOL1d,
		},
		"bnbusdt": map[string]string{
			"1m": config.TelegramTokenBNB1m,
			// "5m":  config.TelegramTokenBNB5m,
			// "15m": config.TelegramTokenBNB15m,
			// "30m": config.TelegramTokenBNB30m,
			// "1h":  config.TelegramTokenBNB1h,
			// "4h": config.TelegramTokenBNB4h,
			// "1d":  config.TelegramTokenBNB1d,
		},
		"adausdt": map[string]string{
			"1m": config.TelegramTokenADA1m,
			// "5m":   config.TelegramTokenADA5m,
			// "15m": config.TelegramTokenADA15m,
			// "30m": config.TelegramTokenADA30m,
			// "1h":  config.TelegramTokenADA1h,
			// "4h": config.TelegramTokenADA4h,
			// "1d":  config.TelegramTokenADA1d,
		},
		"xrpusdt": map[string]string{
			"1m": config.TelegramTokenXRP1m,
			// "5m":   config.TelegramTokenXRP5m,
			// "15m": config.TelegramTokenXRP15m,
			// "30m":  config.TelegramTokenXRP30m,
			// "1h":  config.TelegramTokenXRP1h,
			// "4h":config.TelegramTokenXRP4h,
			// "1d": config.TelegramTokenXRP1d,
		},
	}
}
