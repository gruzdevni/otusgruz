package deliveryhttp

type DeliveryStatus struct {
	OrderID        string `json:"order_id"`
	DeliveryStatus string `json:"delivery_status"`
}

type DefaultResponse struct {
	Status int `json:"status"`
	Errors any `json:"errors"`
}

type ReserveDelivery struct {
	OrderID string `json:"order_id"`
	SlotID  string `json:"slot_id"`
}
