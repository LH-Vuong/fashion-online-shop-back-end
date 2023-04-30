package zalopay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zpmep/hmacutil"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"online_fashion_shop/api/model/order"
	"online_fashion_shop/api/model/payment"
	"strconv"
	"time"
)

type Processor interface {
	InitPayment(info *order.OrderInfo) error
	GetPaymentStatus(paymentId string) (payment.Status, error)
}

type ZaloPaymentProcessor struct {
	appid string
	key1  string
	key2  string
}

func NewZaloPayProcessor(appId string, key1 string, key2 string) Processor {
	return &ZaloPaymentProcessor{
		appid: appId,
		key1:  key1,
		key2:  key2,
	}
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
	info.PaymentInfo.PaymentId = uuid.New().String()
	info.PaymentInfo.PaymentAt = data["app_time"].(int64)
	info.PaymentInfo.PaymentAt = data["server_time"].(int64)
	info.PaymentInfo.Status = payment.StatusApproved
	info.PaymentInfo.ReceivedAmount = data["amount"].(int64)
	info.PaymentInfo.UpdatedAt = time.Now().UnixMilli()
	return nil
}

type object map[string]interface{}

var emp map[string]interface{}

type EmbedData struct {
	RedirectUrl string `json:"redirecturl"`
}

func (processor *ZaloPaymentProcessor) InitPayment(info *order.OrderInfo) error {

	rand.Seed(time.Now().UnixNano())
	transID := rand.Intn(1000000) // Generate random trans id
	embedData, _ := json.Marshal(EmbedData{RedirectUrl: "https://docs.zalopay.vn/result"})
	items, _ := json.Marshal(info.Items)
	amount := strconv.FormatInt(info.TotalPrice, 10)
	// request data
	params := make(url.Values)
	params.Add("app_id", processor.appid)
	params.Add("amount", amount)
	params.Add("app_user", "admin")
	params.Add("embed_data", string(embedData))
	params.Add("item", string(items))
	params.Add("description", "Payment for the order #"+strconv.Itoa(transID))
	params.Add("bank_code", "zalopayapp")

	now := time.Now()
	params.Add("app_time", strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)) // miliseconds

	params.Add("app_trans_id", fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID)) // translation missing: vi.docs.shared.sample_code.comments.app_trans_id

	// appid|app_trans_id|appuser|amount|apptime|embeddata|item
	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, processor.key1, data))

	// Content-Type: application/x-www-form-urlencoded
	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)

	// parse response
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var result payment.ZaloPayApiResult
	println("BODY")
	println(string(body))
	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}

	switch result.ReturnCode {
	case 1:
		info.PaymentInfo.Status = payment.StatusPending
		info.PaymentInfo.PaymentId = fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID)
	default:
		info.PaymentInfo.Status = payment.StatusError
	}
	info.PaymentInfo.PaymentAt = time.Now().UnixMilli()

	return nil
}

func (processor *ZaloPaymentProcessor) GetPaymentStatus(paymentId string) (payment.Status, error) {
	data := fmt.Sprintf("%v|%s|%s", processor.appid, paymentId, processor.key1) // appid|apptransid|key1
	println(data)
	params := map[string]interface{}{
		"app_id":       processor.appid,
		"app_trans_id": paymentId,
		"mac":          hmacutil.HexStringEncode(hmacutil.SHA256, processor.key1, data),
	}

	jsonStr, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var result payment.ZaloPayApiResult
	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.ReturnCode != 1 {
		return payment.StatusError, err
	}

	return payment.StatusApproved, nil

}
