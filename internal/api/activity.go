package api

type Activity struct {
	UserId    string  `json:"user_id"`
	Type      string  `json:"type"`
	ProductId string  `json:"product_id"`
	Quantity  *int    `json:"quantity"`
	Price     *string `json:"price"`
}
