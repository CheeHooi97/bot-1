package constant

func GetThreadIdMap() map[string]any {
	return map[string]any{
		"btcusdt": map[string]int64{
			"1m":  12,
			"5m":  72,
			"15m": 74,
			"1h":  76,
			"4h":  79,
		},
		"ethusdt": map[string]int64{
			"1m":  42,
			"5m":  51,
			"15m": 57,
			"1h":  63,
			"4h":  66,
		},
		"solusdt": map[string]int64{
			"1m":  86,
			"5m":  92,
			"15m": 95,
			// "1h": 0,
			// "4h": 0,
			// "1d": 0,
		},
		"bnbusdt": map[string]int64{
			"1m": 88,
			// "5m":  0,
			"15m": 99,
			// "1h":  0,
			// "4h":0,
			// "1d": 0,
		},
		"adausdt": map[string]int64{
			"1m":  5,
			"5m":  101,
			"15m": 104,
			// "1h": 0,
			// "4h": 0,
			// "1d": 0,
		},
		"xrpusdt": map[string]int64{
			"1m": 90,
			// "5m":   0,
			"15m": 117,
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
