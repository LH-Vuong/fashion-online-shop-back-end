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
	ServiceTypeId  string
}

func (service *GHNService) CalculateFee(districtId string, wardId string, codValue int) (int, error) {
	url := "https://dev-online-gateway.ghn.vn/shiip/public-api/v2/shipping-order/fee"
	method := "POST"

	payload := map[string]interface{}{
		"from_district_id": service.ShopDistrictId,
		"from_ward_code":   service.ShopWardId,
		"service_id":       nil,
		"service_type_id":  service.ServiceTypeId,
		"to_district_id":   districtId,
		"to_ward_code":     wardId,
		"height":           50,
		"length":           20,
		"weight":           200,
		"width":            20,
		"cod_value":        codValue,
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
	{
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("ShopId", service.ShopId)
		req.Header.Add("Token", service.Token)
	}

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
