package payment

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zpmep/hmacutil"
	"io"
	"net/http"
	"net/url"
	"online_fashion_shop/api/model/order"
	"strconv"
	"time"
)

type Processor interface {
	InitPayment(info *order.OrderInfo) error
	GetPaymentStatus(paymentId string) (order.Status, error)
}

type ZaloPaymentProcessor struct {
	appid string
	key1  string
	key2  string
}

func IsValidCallback(requestMac string, dataStr string, key2 string) error {
	mac := hmacutil.HexStringEncode(hmacutil.SHA256, key2, dataStr)

	// kiểm tra callback hợp lệ (đến từ ZaloPay server)
	if mac != requestMac {
		// callback không hợp lệ
		return fmt.Errorf("invalid callback")
	}
	return nil
}

func HandleCallback(info *order.OrderInfo, data map[string]any) error {
	info.PaymentInfo.PaymentAt = data["app_time"].(int64)
	info.PaymentInfo.PaymentAt = data["server_time"].(int64)
	info.PaymentInfo.Status = order.StatusApproved
	info.PaymentInfo.ReceivedAmount = data["amount"].(int64)
	info.PaymentInfo.UpdatedAt = time.Now().UnixMilli()
	return nil
}

func (processor *ZaloPaymentProcessor) InitPayment(info *order.OrderInfo) error {

	embeddata, _ := json.Marshal(nil)
	items := fmt.Sprintf("[{\"order_id\":\"%s\"}]", info.Id)

	// request data
	params := make(url.Values)
	params.Add("appid", processor.appid)
	params.Add("amount", string(info.TotalPrice))
	params.Add("appuser", "demo")
	params.Add("embeddata", string(embeddata))
	params.Add("item", items)
	params.Add("description", "ZaloPay QR Merchant")
	params.Add("bankcode", "zalopayapp")

	now := time.Now()
	params.Add("apptime", strconv.FormatInt(now.UnixMilli(), 10)) // miliseconds

	transid := uuid.New().String()                                                                                 // unique id
	params.Add("apptransid", fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transid)) // mã giao dich có định dạng yyMMdd_xxxx

	// appid|apptransid|appuser|amount|apptime|embeddata|item
	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v",
		params.Get("appid"),
		params.Get("apptransid"),
		params.Get("appuser"),
		params.Get("amount"),
		params.Get("apptime"), params.Get("embeddata"),
		params.Get("item"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, processor.key1, data))

	// Content-Type: application/x-www-form-urlencoded
	res, err := http.PostForm("https://sandbox.zalopay.com.vn/v001/tpe/createorder", params)

	// parse response
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var result order.ZaloPayApiResult

	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}

	switch result.ReturnCode {
	case 1:
		info.PaymentInfo.Status = order.StatusPending
	default:
		info.PaymentInfo.Status = order.StatusError
	}
	info.PaymentInfo.PaymentAt = time.Now().UnixMilli()

	return nil
}

func (processor *ZaloPaymentProcessor) GetPaymentStatus(paymentId string) (order.Status, error) {
	params := make(url.Values)
	params.Add("appid", processor.appid)
	params.Add("apptransid", paymentId)

	data := fmt.Sprintf("%s|%s|%s", processor.appid, params.Get("apptransid"), processor.key1) // appid|apptransid|key1
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, processor.key1, data))

	res, err := http.Get("https://sandbox.zalopay.com.vn/v001/tpe/getstatusbyapptransid?" + params.Encode())

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var result order.ZaloPayApiResult

	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.ReturnCode != 1 {
		return order.StatusError, err
	}

	return order.StatusApproved, nil

}
