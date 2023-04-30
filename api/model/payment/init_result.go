package payment

type ZaloPayApiResult struct {
	ZPTransToken  string `json:"zp_trans_token"`
	OrderUrl      string `json:"order_url"`
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
}
