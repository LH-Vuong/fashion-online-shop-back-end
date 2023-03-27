package container

import (
	"go.uber.org/dig"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/api/service"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	container.Provide(provideMockProductRepository)
	container.Provide(provideProductService)
	return container
}

func provideMockProductRepository() repository.ProductRepository {
	return repository.NewMockProductRepository()
}
func provideProductService(repository repository.ProductRepository) service.ProductService {
	return service.NewProductServiceImpl(repository)
}
