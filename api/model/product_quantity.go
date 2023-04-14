package model

type ProductQuantity struct {
	Id       string `bson:"_id,omitempty"json:"id"`
	Color    string `bson:"color" json:"color"`
	Size     string `bson:"size" json:"size"`
	Quantity int    `bson:"quantity"json:"quantity"`
	DetailId string `bson:"detail_id" json:"detail_id"`
}
