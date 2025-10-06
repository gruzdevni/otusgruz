package goodshttp

type ReserveStatus struct {
	OrderID       string `json:"order_id"`
	ReserveStatus string `json:"reserve_status"`
}

type DefaultResponse struct {
	Status int `json:"status"`
	Errors any `json:"errors"`
}

type ReserveGoods struct {
	OrderID string  `json:"order_id"`
	Goods   []Goods `json:"goods"`
}

type Goods struct {
	Nomenclature string `json:"nomenclature"`
	Quantity     int    `json:"quantity"`
}
