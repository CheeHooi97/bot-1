package constant

func GetThreadIdMap() map[string]any {
	return map[string]any{
		"btcusdt": map[string]int64{
			"1m":  0,
			"5m":  0,
			"15m": 0,
			"1h":  0,
			"4h":  0,
		},
		"ethusdt": map[string]int64{
			"1m":  0,
			"5m":  0,
			"15m": 0,
			"1h":  0,
			"4h":  0,
		},
		"solusdt": map[string]int64{
			"1m":  0,
			"5m":  0,
			"15m": 0,
			// "1h": 0,
			// "4h": 0,
			// "1d": 0,
		},
		"bnbusdt": map[string]int64{
			"1m": 0,
			// "5m":  0,
			"15m": 0,
			// "1h":  0,
			// "4h":0,
			// "1d": 0,
		},
		"adausdt": map[string]int64{
			"1m":  5,
			"5m":  0,
			"15m": 0,
			// "1h": 0,
			// "4h": 0,
			// "1d": 0,
		},
		"xrpusdt": map[string]int64{
			"1m": 0,
			// "5m":   0,
			"15m": 0,
			// "1h": 0,
			// "4h":0,
			// "1d": 0,
		},
		"hypeusdt": map[string]int64{
			"1m": 0,
		},
		"dogeusdt": map[string]int64{
			"1m": 0,
		},
	}
}
