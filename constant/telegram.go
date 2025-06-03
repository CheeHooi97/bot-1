package constant

func GetThreadIdMap() map[string]any {
	return map[string]any{
		"btcusdt": map[string]string{
			"1m":  "",
			"5m":  "",
			"15m": "",
			"1h":  "",
			"4h":  "",
		},
		"ethusdt": map[string]string{
			"1m":  "",
			"5m":  "",
			"15m": "",
			"1h":  "",
			"4h":  "",
		},
		"solusdt": map[string]string{
			"1m":  "",
			"5m":  "",
			"15m": "",
			// "1h": "",
			// "4h": "",
			// "1d": "",
		},
		"bnbusdt": map[string]string{
			"1m": "",
			// "5m":  "",
			"15m": "",
			// "1h":  "",
			// "4h":"",
			// "1d": "",
		},
		"adausdt": map[string]string{
			"1m":  "5",
			"5m":  "",
			"15m": "",
			// "1h": "",
			// "4h": "",
			// "1d": "",
		},
		"xrpusdt": map[string]string{
			"1m": "",
			// "5m":   "",
			"15m": "",
			// "1h": "",
			// "4h":"",
			// "1d": "",
		},
		"hypeusdt": map[string]string{
			"1m": "",
		},
		"dogeusdt": map[string]string{
			"1m": "",
		},
	}
}
