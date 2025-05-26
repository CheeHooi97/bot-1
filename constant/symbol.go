package constant

import "bot-1/config"

// var TokenMap = map[string][]string{
// 	"btcusdt": {
// 		config.TelegramTokenBTC1m,  //0
// 		config.TelegramTokenBTC5m,  //1
// 		config.TelegramTokenBTC15m, //2
// 		config.TelegramTokenBTC30m, //3
// 		config.TelegramTokenBTC1h,  //4
// 		config.TelegramTokenBTC4h,  //5
// 		// config.TelegramTokenBTC1d,  //6
// 	},
// 	"ethusdt": {
// 		config.TelegramTokenETH1m,  //0
// 		config.TelegramTokenETH5m,  //1
// 		config.TelegramTokenETH15m, //2
// 		config.TelegramTokenETH30m, //3
// 		config.TelegramTokenETH1h,  //4
// 		config.TelegramTokenETH4h,  //5
// 		// config.TelegramTokenETH1d,  //6
// 	},
// 	"solusdt": {
// 		config.TelegramTokenSOL1m, //0
// 		// config.TelegramTokenSOL5m, //1
// 		// config.TelegramTokenSOL15m, //2
// 		// config.TelegramTokenSOL30m, //3
// 		// config.TelegramTokenSOL1h, //4
// 		// config.TelegramTokenSOL4h, //5
// 		// config.TelegramTokenSOL1d, //06
// 	},
// 	"bnbusdt": {
// 		config.TelegramTokenBNB1m, //0
// 		// config.TelegramTokenBNB5m, //1
// 		// config.TelegramTokenBNB15m, //2
// 		// config.TelegramTokenBNB30m, //3
// 		// config.TelegramTokenBNB1h, //4
// 		// config.TelegramTokenBNB4h, //5
// 		// config.TelegramTokenBNB1d, //06
// 	},
// 	"adausdt": {
// 		config.TelegramTokenADA1m, //0
// 		// config.TelegramTokenADA5m, //1
// 		// config.TelegramTokenADA15m, //2
// 		// config.TelegramTokenADA30m, //3
// 		// config.TelegramTokenADA1h, //4
// 		// config.TelegramTokenADA4h, //5
// 		// config.TelegramTokenADA1d, //06
// 	},
// 	"xrpusdt": {
// 		config.TelegramTokenXRP1m, //0
// 		// config.TelegramTokenXRP5m, //1
// 		// config.TelegramTokenXRP15m, //2
// 		// config.TelegramTokenXRP30m, //3
// 		// config.TelegramTokenXRP1h, //4
// 		// config.TelegramTokenXRP4h, //5
// 		// config.TelegramTokenXRP1d, //06
// 	},
// }

var TokenMap = map[string]any{
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
