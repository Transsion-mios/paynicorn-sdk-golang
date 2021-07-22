package paynicorn

import (
	"fmt"
	"github.com/Transsion-mios/paynicorn-sdk-golang"
	"github.com/gin-gonic/gin"
)

func main() {
	var appkey = "PUT_YOUR_APPKEY_HERE"
	var merchantkey = "PUT_YOUR_MERCHANT_SECRET_HERE"


	//raise a payment request to PAYNICORN
	cashierRequest := paynicorn.RaiseCashierRequest{}
	cashierRequest.OrderId="PUT_YOUR_ORDER_ID_HERE"
	cashierRequest.CountryCode="NG"
	cashierRequest.Currency="NGN"
	cashierRequest.Amount="10"
	cashierRequest.CpFrontPage="PUT_YOUR_WEB_REDIRECT_URL_HERE"
	cashierRequest.OrderDescription="TEST GOODS NAME"
	cashierresponse,_ := paynicorn.RaiseCashierPayment(appkey,merchantkey,cashierRequest)
	fmt.Println(cashierresponse)


	//query a payment status from PAYNICORN
	request := paynicorn.QueryTransactionRequest{}
	request.OrderId=cashierRequest.OrderId
	request.TxnType=paynicorn.PAYMENT
	response := paynicorn.QueryPaymentStatus(appkey,merchantkey,request)
	fmt.Println(response)


	//receive a payment status postback from PAYNICORN
	r := gin.Default()
	r.POST("/postback", func(context *gin.Context) {
		postbackresponse,_ := paynicorn.Postback(context,merchantkey)
		fmt.Println(postbackresponse)

	})
	r.Run(":80")


}
