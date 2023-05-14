package rating

type AvrRate struct {
	ProductId string  `bson:"_id"`
	AvrRate   float64 `bson:"avr_rate"`
}
