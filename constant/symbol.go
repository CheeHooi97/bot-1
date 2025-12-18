package constant

import "bot-1/config"

func GetTokenMap() map[string]any {
	return map[string]any{
		"btcusdt": map[string]string{
			"4h": config.BTCFuture4h,
			"1h": config.BTCFuture1h,
		},
		"ethusdt": map[string]string{
			"4h": config.ETHFuture4h,
			"1h": config.ETHFuture1h,
		},
		"adausdt": map[string]string{
			"4h": config.ADAFuture4h,
		},
		"solusdt": map[string]string{
			"4h": config.SOLFuture4h,
		},
		"bnbusdt": map[string]string{
			"4h": config.BNBFuture4h,
		},
	}
}
