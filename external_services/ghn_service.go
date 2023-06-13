package external_services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GHNService struct {
	Token          string
	ShopId         string
	ShopDistrictId string
	ShopWardId     string
	ServiceTypeId  int
}

func (service *GHNService) CalculateFee(districtId int, wardCode string) (int, error) {
	url := "https://online-gateway.ghn.vn/shiip/public-api/v2/shipping-order/fee"
	method := "POST"

	payload := map[string]interface{}{
		"from_district_id": 1454,
		"from_ward_code":   "21211",
		"service_id":       nil,
		"service_type_id":  2,
		"to_district_id":   districtId,
		"to_ward_code":     wardCode,
		"height":           50,
		"length":           20,
		"weight":           200,
		"width":            20,
		"cod_value":        nil,
		"coupon":           nil,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payloadBytes))
	defer req.Body.Close()
	if err != nil {
		return 0, err
	}
	//set header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ShopId", service.ShopId)
	req.Header.Add("Token", service.Token)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	var responseMap struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Total int `json:"total"`
		} `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&responseMap)
	if err != nil {
		return 0, err
	}
	if responseMap.Code != 200 {
		return 0, fmt.Errorf("GHN ERROR_RESPONSE, MESSAGE(%s)", responseMap.Message)
	}
	return responseMap.Data.Total, nil
}
