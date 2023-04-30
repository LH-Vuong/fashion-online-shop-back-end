package container

import (
	"log"
	"online_fashion_shop/api/repository"
	userrepo "online_fashion_shop/api/repository/user"
	"online_fashion_shop/api/service"
	"online_fashion_shop/initializers"
	"online_fashion_shop/initializers/zalopay"

	"go.uber.org/dig"
)

var container = dig.New()

func init() {
	container.Provide(provideMongoDbClient)
	container.Provide(provideProductRepositoryImpl)
	container.Provide(provideProductRatingRepositoryImpl)
	container.Provide(providePhotoRepositoryImpl)
	container.Provide(provideUserRepositoryImpl)
	container.Provide(provideProductQuantityRepositoryImpl)
	container.Provide(provideCartRepositoryImpl)
	container.Provide(provideProductServiceImpl)
	container.Provide(provideCardServiceImpl)
	container.Provide(providePhotoServiceImpl)
	container.Provide(provideCouponService)
	container.Provide(provideCouponRepositoryImpl)
	container.Provide(provideUserServiceImpl)
	container.Provide(provideOrderService)
	container.Provide(provideOrderRepositoryImpl)
	container.Provide(provideZaloPayProcessor)

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

func providePhotoServiceImpl(photoRepo repository.ProductPhotoRepository) service.PhotoService {
	return service.NewPhotoServiceImpl(photoRepo)
}

func providePhotoRepositoryImpl(client initializers.Client) repository.ProductPhotoRepository {
	productPhotoCollection := client.Database("fashion_shop").Collection("product_photo")
	return repository.NewProductPhotoRepository(productPhotoCollection)
}
func provideProductServiceImpl(detailRepo repository.ProductDetailRepository,
	photoService service.PhotoService,
	ratingRepo repository.ProductRatingRepository) service.ProductService {
	return service.NewProductServiceImpl(detailRepo, ratingRepo, photoService)
}

func provideMongoDbClient() initializers.Client {
	config, err := initializers.LoadConfig("../")

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

func provideUserServiceImpl(userRepo userrepo.UserRepository) service.UserService {
	return service.NewUserServiceImpl(userRepo)
}

func provideCartRepositoryImpl(cl initializers.Client) repository.CartRepository {
	cartCollection := cl.Database("fashion_shop").Collection("cart")
	return repository.NewCartRepositoryImpl(cartCollection)
}

func provideCouponRepositoryImpl(cl initializers.Client) repository.CouponRepository {
	couponCollection := cl.Database("fashion_shop").Collection("coupon")
	return repository.NewCouponRepositoryImpl(couponCollection)
}

func provideOrderRepositoryImpl(cl initializers.Client) repository.OrderRepository {
	orderInfo := cl.Database("fashion_shop").Collection("order")
	return repository.NewOrderRepositoryImpl(orderInfo)
}

func provideOrderService(orderRepo repository.OrderRepository,
	cartService service.CartService,
	couponService service.CouponService,
	processor zalopay.Processor,
) service.OrderService {
	return service.NewOrderServiceImpl(couponService, cartService, orderRepo, processor)
}

func provideCouponService(couponRepo repository.CouponRepository,
) service.CouponService {
	return service.NewCouponService(couponRepo)
}

func provideZaloPayProcessor() zalopay.Processor {
	config, err := initializers.LoadConfig("../")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
		panic(err)
	}
	return zalopay.NewZaloPayProcessor(config.ZaloAppId, config.ZaloKey1, config.ZaloKey2)

}

func provideUserRepositoryImpl(cl initializers.Client) userrepo.UserRepository {
	userCollection := cl.Database("fashion_shop").Collection("user")
	userVerifyCollection := cl.Database("fashion_shop").Collection("user_verify")
	userWishlistCollection := cl.Database("fashion_shop").Collection("user_wishlist")
	userAddressCollection := cl.Database("fashion_shop").Collection("user_address")
	return userrepo.NewUserRepositoryImpl(userCollection, userVerifyCollection, userWishlistCollection, userAddressCollection)
}
