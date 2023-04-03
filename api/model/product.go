package model

type Product struct {
	Id              string         `bson:"_id"`
	Name            string         `bson:"name"`
	Tags            []string       `bson:"tags"`
	Types           []string       `bson:"types"`
	Brand           string         `bson:"brand"`
	DiscountAmount  float64        `json:"discount_amount"`
	DiscountPercent float64        `json:"discount_percent"`
	Gender          string         `bson:"gender"`
	Price           int64          `json:"price"`
	Description     string         `bson:"description"`
	Photos          []ProductPhoto `bson:"photos"`
	AvrRate         float64        `json:"avr_rate"`
}
