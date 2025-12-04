package constant

import "dca-bot/config"

func GetTokenMap() map[string]any {
	return map[string]any{
		"btcusdt": map[float32]string{
			2:   config.BTC2,
			1.5: config.BTC1x5,
			1:   config.BTC1,
		},
	}
}
