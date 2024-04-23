package kraken

type TickerResponse struct {
	Error  []interface{}    `json:"error"`
	Result map[string]*Data `json:"result"`
}

type Data struct {
	A []string `json:"a"`
	B []string `json:"b"`
	C []string `json:"c"`
	V []string `json:"v"`
	P []string `json:"p"`
	T []int    `json:"t"`
	L []string `json:"l"`
	H []string `json:"h"`
	O string   `json:"o"`
}
