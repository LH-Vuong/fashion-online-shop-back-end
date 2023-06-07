package product

type Product struct {
	Id                string             `bson:"_id,omitempty" json:"id"`
	Name              string             `bson:"name,omitempty" json:"name"`
	Tags              []string           `bson:"tags,omitempty" json:"tags"`
	Types             []string           `bson:"types,omitempty" json:"types"`
	Brand             string             `bson:"brand,omitempty" json:"brand"`
	DiscountAmount    float64            `bson:"discount_amount,omitempty" json:"-"`
	DiscountPercent   float64            `bson:"discount_percent,omitempty" json:"discount_percent"`
	Gender            string             `bson:"gender,omitempty" json:"gender"`
	Price             int64              `bson:"price,omitempty" json:"price"`
	Description       string             `bson:"description,omitempty" json:"description"`
	Photos            []string           `bson:"-" json:"photos"`
	ProductQuantities []*ProductQuantity `json:"product_quantities" bson:"-"`
	AvrRate           int                `bson:"-" json:"avr_rate"`
	CreatedAt         int64              `bson:"create_at,omitempty" json:"-"`
	CreateBy          string             `bson:"create_by,omitempty" json:"-"`
	UpdatedAt         int64              `bson:"updated_at,omitempty" json:"-"`
	UpdatedBy         string             `bson:"updated_by,omitempty" json:"-"`
}
