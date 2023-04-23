package test

import (
	"log"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/api/service"
	"online_fashion_shop/initializers"

	"go.uber.org/dig"
)

var container = dig.New()

func init() {
	container.Provide(provideMongoDbClient)
	container.Provide(provideProductRepositoryImpl)
	container.Provide(provideProductRatingRepositoryImpl)
	container.Provide(providePhotoRepositoryImpl)
	container.Provide(provideProductServiceImpl)
	container.Provide(provideProductQuantityRepositoryImpl)
	container.Provide(provideCartRepositoryImpl)
	container.Provide(provideCardServiceImpl)
}

func BuildContainer() *dig.Container {
	return container
}

func Inject[A any](dependency *A) error {
	return container.Invoke(func(value A) {
		*dependency = value
	})
}

func provideProductRatingRepositoryImpl(client initializers.Client) repository.ProductRatingRepository {
	ratingCollection := client.Database("fashion_shop").Collection("product_rating")
	return repository.NewProductRatingRepositoryImpl(ratingCollection)
}

func provideProductRepositoryImpl(mongoClient initializers.Client) repository.ProductDetailRepository {
	productCollection := mongoClient.Database("fashion_shop").Collection("product")
	return repository.NewProductRepositoryImpl(productCollection)
}

func providePhotoRepositoryImpl(client initializers.Client) repository.ProductPhotoRepository {
	productPhotoCollection := client.Database("fashion_shop").Collection("product_photo")
	return repository.NewProductPhotoRepository(productPhotoCollection)
}
func provideProductServiceImpl(detailRepo repository.ProductDetailRepository,
	photoRepo repository.ProductPhotoRepository,
	ratingRepo repository.ProductRatingRepository) service.ProductService {
	return service.NewProductServiceImpl(detailRepo, ratingRepo, photoRepo)
}

func provideMongoDbClient() initializers.Client {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	client, err := initializers.NewClient(config.MongoUrl)
	if err != nil {
		panic(err)
	}
	return client
}

func provideProductQuantityRepositoryImpl(cl initializers.Client) repository.ProductQuantityRepository {
	quantityCollection := cl.Database("fashion_shop").Collection("product_quantity")
	return repository.NewProductQuantityRepositoryImpl(quantityCollection)
}

func provideCardServiceImpl(cartRepo repository.CartRepository,
	quantityRepo repository.ProductQuantityRepository,
	detailRepo repository.ProductDetailRepository) service.CartService {
	return service.NewCartServiceImpl(cartRepo, quantityRepo, detailRepo)
}

func provideCartRepositoryImpl(cl initializers.Client) repository.CartRepository {
	cartCollection := cl.Database("fashion_shop").Collection("cart")
	return repository.NewCartRepositoryImpl(cartCollection)
}
