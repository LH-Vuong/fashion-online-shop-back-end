package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"online_fashion_shop/api/controller"
	"online_fashion_shop/api/service"
)

func InitDeliveryRouter(s *gin.Engine, container *dig.Container) {
	err := container.Invoke(func(deliveryService service.DeliveryService) {
		deliveryCtrl := controller.DeliveryController{DeliveryService: deliveryService}
		s.GET("api/delivery/cal_fee/:address_id", deliveryCtrl.CalculateFee)
	})
	if err != nil {
		panic(err)
	}

}
