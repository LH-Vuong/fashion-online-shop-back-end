package product

type ProductQuantity struct {
	Id       string `bson:"_id,omitempty" json:"id"`
	Color    string `bson:"color,omitempty" json:"color"`
	Size     string `bson:"size,omitempty" json:"size"`
	Quantity int    `bson:"quantity,omitempty"json:"quantity"`
	DetailId string `bson:"detail_id,omitempty" json:"detail_id"`
}
