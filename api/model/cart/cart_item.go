package cart

import "online_fashion_shop/api/model/product"

type CartItem struct {
	CustomerId    string          `bson:"customer_id,omitempty" json:"customer_id"`
	InventoryId   string          `bson:"product_id" json:"inventory_id"`
	Quantity      int             `bson:"quantity" json:"quantity"`
	Color         string          `bson:"color" json:"color"`
	Size          string          `bson:"size" json:"size"`
	ProductDetail product.Product `bson:"product_detail" json:"product_detail"`
	CreatedAt     int64           `bson:"create_at,omitempty" json:"-"`
	CreateBy      string          `bson:"create_by,omitempty" json:"-"`
	UpdatedAt     int64           `bson:"updated_at,omitempty" json:"-"`
	UpdatedBy     string          `bson:"updated_by,omitempty" json:"-"`
}
