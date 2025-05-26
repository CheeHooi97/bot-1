package constant

var SymbolPrecisionMap = map[string][]int{
	"btcusdt": {2, 4}, // pricePrecision = 2, amountPrecision = 4
	"ethusdt": {2, 2},
	"adausdt": {4, 2},
	"bnbusdt": {2, 2},
	"solusdt": {2, 2},
	"xrpusdt": {2, 2},
}

var StepMap = map[string]float64{
	"btcusdt": 0.001,
	"ethusdt": 0.01,
	"adausdt": 0.1,
	"bnbusdt": 0.1,
	"solusdt": 0.1,
	"xrpusdt": 0.1,
}
