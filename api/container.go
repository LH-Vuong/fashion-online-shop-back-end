package container

import (
	"go.uber.org/dig"
	"online_fashion_shop/api/dbs"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/api/service"
	"online_fashion_shop/configs"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	container.Provide(provideMongoDbClient)
	container.Provide(provideProductRepositoryImpl)
	container.Provide(provideProductRatingRepositoryImpl)
	container.Provide(providePhotoRepositoryImpl)
	container.Provide(provideProductService)
	return container
}

func provideProductRatingRepositoryImpl(client dbs.Client) repository.ProductRatingRepository {
	ratingCollection := client.Database("fashion_shop").Collection("product_rating")
	return repository.NewProductRatingRepositoryImpl(ratingCollection)
}

func provideProductRepositoryImpl(mongoClient dbs.Client) repository.ProductDetailRepository {
	productCollection := mongoClient.Database("fashion_shop").Collection("product")
	return repository.NewProductRepositoryImpl(productCollection)
}

func providePhotoRepositoryImpl(client dbs.Client) repository.ProductPhotoRepository {
	productPhotoCollection := client.Database("fashion_shop").Collection("product_photo")
	return repository.NewProductPhotoRepository(productPhotoCollection)
}
func provideProductService(detailRepo repository.ProductDetailRepository,
	photoRepo repository.ProductPhotoRepository,
	ratingRepo repository.ProductRatingRepository) service.ProductService {
	return service.NewProductServiceImpl(detailRepo, ratingRepo, photoRepo)
}

func provideMongoDbClient() dbs.Client {
	client, err1 := dbs.NewClient(configs.GetString("mongo_db.url"))
	if err1 != nil {
		panic(err1)
	}
	return client
}
