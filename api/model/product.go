package model

type Product struct {
	Id              string   `bson:"_id" json:"id"`
	Name            string   `bson:"name" json:"name"`
	Tags            []string `bson:"tags" json:"tags"`
	Types           []string `bson:"types" json:"types"`
	Brand           string   `bson:"brand" json:"brand"`
	DiscountAmount  float64  `bson:"discount_amount" json:"discount_amount"`
	DiscountPercent float64  `bson:"discount_percent" json:"discount_percent"`
	Gender          string   `bson:"gender" json:"gender"`
	Price           int64    `bson:"price" json:"price"`
	Description     string   `bson:"description" json:"description"`
	Photos          []string `bson:"-" json:"photos"`
	AvrRate         float64  `bson:"-" json:"avr_rate"`
	CreatedAt       int64    `bson:"create_at,omitempty" json:"-"`
	CreateBy        string   `bson:"create_by,omitempty" json:"-"`
	UpdatedAt       int64    `bson:"updated_at,omitempty" json:"-"`
	UpdatedBy       string   `bson:"updated_by,omitempty" json:"-"`
}
