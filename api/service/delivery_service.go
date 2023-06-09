package service

import (
	"context"
	"fmt"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/external_services"
)

type DeliveryService interface {
	CalculateFeeByCustomerAddress(addressId string) (int, error)
	GetDeliveryAddress(addressId string) (string, error)
}

type DeliveryServiceImpl struct {
	gnhService *external_services.GHNService
	userRepo   repository.UserRepository
}

func (deliveryService *DeliveryServiceImpl) GetDeliveryAddress(addressId string) (string, error) {
	address, err := deliveryService.userRepo.GetUserAddressById(context.TODO(), addressId)
	if err != nil {
		return "", fmt.Errorf("encoutered error(%s) while trying to retrive address info's(%s)", err.Error(), addressId)
	}
	return address.Address, nil
}

func NewDeliveryServiceImpl(gnhService *external_services.GHNService,
	userRepo repository.UserRepository) DeliveryService {
	return &DeliveryServiceImpl{
		gnhService: gnhService,
		userRepo:   userRepo,
	}
}

func (deliveryService *DeliveryServiceImpl) CalculateFeeByCustomerAddress(addressId string) (int, error) {
	address, err := deliveryService.userRepo.GetUserAddressById(context.TODO(), addressId)
	if err != nil {
		return 0, fmt.Errorf("encoutered error(%s) while trying to retrive address info's(%s)", err.Error(), addressId)
	}
	return deliveryService.gnhService.CalculateFee(address.DistrictId, address.WardCode)
}
