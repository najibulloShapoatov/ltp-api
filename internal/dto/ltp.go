package dto

type LTPResponse struct {
	Ltp []*LTP `json:"ltp"`
}

type LTP struct {
	Pair   string `json:"pair"`
	Amount string `json:"amount"`
}
