package model

type ProductPhoto struct {
	MainPhoto string   `bson:"main_photo"`
	SubPhotos []string `bson:"sub_photos"`
	Color     string   `bson:"color"`
	ProductId string   `bson:"product_id"`
}
