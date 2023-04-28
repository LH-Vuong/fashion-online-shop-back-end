package zalopay

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zpmep/hmacutil"
	"io"
	"log"
	"net/http"
	"net/url"
	"online_fashion_shop/api/model"
	"online_fashion_shop/api/model/payment"
	"strconv"
	"time"
)

type Config struct {
	appid string
	key1  string
	key2  string
}

func HandlerCallBack(requestMac string, dataStr string, conf Config) error {
	mac := hmacutil.HexStringEncode(hmacutil.SHA256, conf.key2, dataStr)

	result := make(map[string]interface{})

	// kiểm tra callback hợp lệ (đến từ ZaloPay server)
	if mac != requestMac {
		// callback không hợp lệ
		result["returncode"] = -1
		result["returnmessage"] = "mac not equal"
	} else {
		// thanh toán thành công
		result["returncode"] = 1
		result["returnmessage"] = "success"

		// merchant cập nhật trạng thái cho đơn hàng

		var dataJSON map[string]interface{}
		json.Unmarshal([]byte(dataStr), &dataJSON)
		log.Println("update order's status = success where apptransid =", dataJSON["apptransid"])
	}
	return nil
}

func InitPayment(info *model.OrderInfo, config Config) error {

	embeddata, _ := json.Marshal(nil)
	items := fmt.Sprintf("[{\"order_id\":\"%s\"}]", info.Id)

	// request data
	params := make(url.Values)
	params.Add("appid", config.appid)
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
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, config.key1, data))

	// Content-Type: application/x-www-form-urlencoded
	res, err := http.PostForm("https://sandbox.zalopay.com.vn/v001/tpe/createorder", params)

	// parse response
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var result payment.ZaloPayApiResult

	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}

	switch result.ReturnCode {
	case 1:
		info.PaymentInfo.Status = payment.StatusPending
	default:
		info.PaymentInfo.Status = payment.StatusError
	}
	info.PaymentInfo.LastUpdateAt = time.Now().UnixMilli()

	return nil
}

func GetPaymentStatus(paymentId string, config Config) (payment.Status, error) {
	params := make(url.Values)
	params.Add("appid", config.appid)
	params.Add("apptransid", paymentId)

	data := fmt.Sprintf("%s|%s|%s", config.appid, params.Get("apptransid"), config.key1) // appid|apptransid|key1
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, config.key1, data))

	res, err := http.Get("https://sandbox.zalopay.com.vn/v001/tpe/getstatusbyapptransid?" + params.Encode())

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
