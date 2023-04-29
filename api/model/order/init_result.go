package order

type ZaloPayApiResult struct {
	ZPTransToken  string `json:"zptranstoken"`
	OrderUrl      string `json:"orderurl"`
	ReturnCode    int    `json:"returncode"`
	ReturnMessage string `json:"returnmessage"`
}

//	"zptranstoken": "190613000002244_order",
//	"orderurl": "https://qcgateway.zalopay.vn/openinapp?order=eyJ6cHRyYW5zdG9rZW4iOiIxOTA2MTMwMDAwMDIyNDRfb3JkZXIiLCJhcHBpZCI6NTUzfQ",
//	"returncode": 1,
//	"returnmessage": ""
