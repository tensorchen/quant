package entity

type Trade struct {
	// Ticker 交易商品代码
	Ticker string `json:"ticker"`
	// Exchange 交易所名称
	Exchange string `json:"exchange"`

	// TimeNow 警报的当前触发时间
	TimeNow string `json:"timenow"`

	SymInfo SymInfo `json:"syminfo"`

	Strategy Strategy `json:"strategy"`
}

type SymInfo struct {
	// Currency 商品的货币代码，例如：EUR、USD
	Currency string `json:"currency"`
}

type Strategy struct {
	PositionSize string `json:"position_size"`
	Order        Order  `json:"order"`

	MarketPosition         string `json:"market_position"`
	MarketPositionSize     string `json:"market_position_size"`
	PrevMarketPosition     string `json:"prev_market_position"`
	PrevMarketPositionSize string `json:"prev_market_position_size"`
}

type Order struct {
	// Action 订单执行操作，例如："buy"，"sell"
	Action string `json:"action"`
	// Contracts 订单合约数量
	Contracts string `json:"contracts"`

	// Price 订单价格
	Price string `json:"price"`

	ID           string `json:"id"`
	Comment      string `json:"comment"`
	AlertMessage string `json:"alert_message"`
}
